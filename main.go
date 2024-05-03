package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	_ "modernc.org/sqlite"
)

const webSocketUrl string = "wss://ws-us2.pusher.com/app/eb1d5f283081a78b932c?protocol=7&client=js&version=7.6.0&flash=false" // Replace this with the actual URL of the WebSocket server.

func startListeningChat(db *sql.DB, streamerName string, chatRoomId string, dataChannel chan<- Data) {
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

	// Send a message to the WebSocket server.
	message := []byte(fmt.Sprintf("{\"event\":\"pusher:subscribe\",\"data\":{\"auth\":\"\",\"channel\":\"chatrooms.%s.v2\"}}", chatRoomId))
	err = wsutil.WriteClientMessage(conn, ws.OpText, message)
	if err != nil {
		log.Fatalf("error sending message: %v", err)
	}

	// Read messages from the server in a separate goroutine.
	go func() {
		for {
			msg, _, err := wsutil.ReadServerData(conn)
			if err != nil {
				log.Printf("error reading message: %v", err)
				return // or handle reconnection logic here
			}

			var event Message
			if err := json.Unmarshal(msg, &event); err != nil {
				fmt.Println("Error unmarshaling event:", err)
				return
			}

			var data Data
			if err := json.Unmarshal([]byte(event.Data), &data); err != nil {
				fmt.Println("Error unmarshaling data:", err)
				return
			}
			if data.Type != "message" {
				continue
			}

			dataChannel <- data

			fmt.Printf("%s:%s:%s \n", streamerName, data.Sender.Username, data.Content)
		}
	}()

	// Block main goroutine without using time.After
	select {}

}

func main() {
	db, err := sql.Open("sqlite", "file:examplev2.db?cache=shared")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
    CREATE TABLE IF NOT EXISTS messages (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        content TEXT NOT NULL,
        username TEXT NOT NULL,
        channel TEXT NOT NULL,
        date TEXT NOT NULL
    );
    `
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
	dataChannel := make(chan Data, 100)
	go WriteToSqliteFromChannel(db, dataChannel)
	for streamerName, chatRoomId := range StreamerChatSocketUrlList {
		go startListeningChat(db, streamerName, chatRoomId, dataChannel)
	}

	// Block main goroutine
	select {}
}

func WriteToDb(db *sql.DB, dataList []Data) {
	if len(dataList) == 0 {
		return
	}

	valueStrings := make([]string, 0, len(dataList))
	valueArgs := make([]interface{}, 0, len(dataList)*4) // Each data has 4 fields

	// Construct the query string with placeholders for each data
	for _, data := range dataList {
		valueStrings = append(valueStrings, "(?, ?, ?, ?)")
		valueArgs = append(valueArgs, data.Content, data.Sender.Username, data.ChatroomID, data.CreatedAt)
	}

	stmt := fmt.Sprintf("INSERT INTO messages (content, username, channel, date) VALUES %s",
		strings.Join(valueStrings, ", "))

	// Execute the query with all parameters
	_, err := db.Exec(stmt, valueArgs...)
	if err != nil {
		log.Print("Failed to insert multiple rows:", err)
	}
}

func WriteToSqliteFromChannel(db *sql.DB, dataChannel <-chan Data) {
	dataBatch := make([]Data, 0, 10)
	batchSize := 10

	for data := range dataChannel {
		dataBatch = append(dataBatch, data)
		if len(dataBatch) >= batchSize {
			fmt.Println("Batch filled Writing to DB")
			WriteToDb(db, dataBatch)
			dataBatch = make([]Data, 0, 10)
		}
	}
}
