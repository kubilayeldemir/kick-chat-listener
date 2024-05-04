package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"log"
	"math/rand"
	_ "modernc.org/sqlite"
	"net/http"
	"strings"
	"time"
)

const batchSize = 10
const webSocketUrl string = "wss://ws-us2.pusher.com/app/eb1d5f283081a78b932c?protocol=7&client=js&version=7.6.0&flash=false"
const chatroomSubcribeCommand string = "{\"event\":\"pusher:subscribe\",\"data\":{\"auth\":\"\",\"channel\":\"chatrooms.%d.v2\"}}"

func main() {
	db, err := sql.Open("sqlite", "file:sqlite-data/database.db?cache=shared")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = CreateTable(db); err != nil {
		return
	}

	dataChannel := make(chan Data, 100)
	go HandleDataChannelInserts(db, dataChannel)

	for channelName, chatRoomId := range ChannelAndChatIdMap {
		go startListeningChat(channelName, chatRoomId, dataChannel)
	}

	// Block main goroutine
	select {}
}

func CreateTable(db *sql.DB) error {
	sqlStmt := `
    CREATE TABLE IF NOT EXISTS messages (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        content TEXT NOT NULL,
        username TEXT NOT NULL,
        channel TEXT NOT NULL,
        date TEXT NOT NULL
    );
    `
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}
	return err
}

func WriteToDb(db *sql.DB, dataList []Data) {
	defer timer("WriteToDb")()
	if len(dataList) == 0 {
		return
	}

	valueStrings := make([]string, 0, len(dataList))
	valueArgs := make([]interface{}, 0, len(dataList)*4) // Each data has 4 fields

	// Construct the query string with placeholders for each data
	for _, data := range dataList {
		valueStrings = append(valueStrings, "(?, ?, ?, ?)")
		valueArgs = append(valueArgs, data.Content, data.Sender.Username, ChatIdAndChannelMap[data.ChatroomID], data.CreatedAt)
	}

	stmt := fmt.Sprintf("INSERT INTO messages (content, username, channel, date) VALUES %s",
		strings.Join(valueStrings, ", "))

	// Execute the query with all parameters
	_, err := db.Exec(stmt, valueArgs...)
	if err != nil {
		log.Print("Failed to insert multiple rows:", err)
	}
}

func HandleDataChannelInserts(db *sql.DB, dataChannel <-chan Data) {
	dataBatch := make([]Data, 0, batchSize)
	start := time.Now()
	for data := range dataChannel {
		dataBatch = append(dataBatch, data)
		if len(dataBatch) >= batchSize {
			color.White("-> Batch filled at %v\n", time.Since(start))
			start = time.Now()
			WriteToDb(db, dataBatch)
			dataBatch = make([]Data, 0, batchSize)
		}
	}
}

func startListeningChat(streamerName string, chatRoomId int, dataChannel chan<- Data) {
	req, err := http.NewRequest(http.MethodGet, webSocketUrl, nil)
	if err != nil {
		log.Fatalf("error creating request: %v", err)
	}

	// Connect to the server.
	conn, _, _, err := ws.Dialer{}.Dial(req.Context(), webSocketUrl)
	if err != nil {
		log.Fatalf("error connecting to websocket: %v", err)
	}
	defer conn.Close()

	// Send a message to the WebSocket server to subscribe to chatroom.
	message := []byte(fmt.Sprintf(chatroomSubcribeCommand, chatRoomId))
	err = wsutil.WriteClientMessage(conn, ws.OpText, message)
	if err != nil {
		log.Fatalf("error sending message: %v", err)
	}

	c := GetRandomColorForLog()

	for {
		msg, _, err := wsutil.ReadServerData(conn)
		if err != nil {
			log.Printf("error reading message: %v", err)
			return
		}

		go UnmarshallAndSendToChannel(streamerName, msg, dataChannel, c)
	}
}

func UnmarshallAndSendToChannel(streamerName string, msgByte []byte, dataChannel chan<- Data, c *color.Color) {
	var event Message
	if err := json.Unmarshal(msgByte, &event); err != nil {
		fmt.Println("Error unmarshaling event:", err)
		return
	}

	var data Data
	if err := json.Unmarshal([]byte(event.Data), &data); err != nil {
		fmt.Println("Error unmarshaling data:", err)
		return
	}
	if data.Type != "message" {
		return
	}

	dataChannel <- data

	c.Printf("%s:%s:%s \n", streamerName, data.Sender.Username, data.Content)
}

func GetRandomColorForLog() *color.Color {
	colors := []*color.Color{
		color.New(color.FgRed),
		color.New(color.FgGreen),
		color.New(color.FgYellow),
		color.New(color.FgCyan),
		color.New(color.FgBlue),
		color.New(color.FgMagenta),
	}
	return colors[rand.Intn(len(colors))]
}

// timer returns a function that prints the name argument and
// the elapsed time between the call to timer and the call to
// the returned function. The returned function is intended to
// be used in a defer statement:
//
//	defer timer("sum")()
func timer(name string) func() {
	start := time.Now()
	return func() {
		color.White("-> %s took %v\n", name, time.Since(start))
	}
}
