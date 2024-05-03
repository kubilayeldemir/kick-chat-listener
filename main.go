package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "file:example.db?cache=shared")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
    CREATE TABLE IF NOT EXISTS messages (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        content TEXT NOT NULL,
        username TEXT NOT NULL
    );
    `
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	// URL of the WebSocket server.
	url := "wss://ws-us2.pusher.com/app/eb1d5f283081a78b932c?protocol=7&client=js&version=7.6.0&flash=false" // Replace this with the actual URL of the WebSocket server.

	// Create a request to establish a WebSocket connection.
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalf("error creating request: %v", err)
	}

	// Connect to the server.
	conn, _, _, err := ws.Dialer{}.Dial(req.Context(), url)
	if err != nil {
		log.Fatalf("error connecting to websocket: %v", err)
	}
	defer conn.Close()

	// Send a message to the WebSocket server.
	message := []byte("{\"event\":\"pusher:subscribe\",\"data\":{\"auth\":\"\",\"channel\":\"chatrooms.25314085.v2\"}}")
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

			//fmt.Printf("Parsed Data: %+v\n", data)
			fmt.Printf("%s:%s \n", data.Sender.Username, data.Content)
			_, err = db.Exec("INSERT INTO messages(content, username) VALUES(?, ?)", data.Content, data.Sender.Username)
			if err != nil {
				print("Error inserting message: %v \n", err)
			}
		}
	}()

	// Block main goroutine without using time.After
	select {}
}
