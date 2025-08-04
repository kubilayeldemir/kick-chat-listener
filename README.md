High performant chat listener that listens multiple channels at the same time, saves every message sent to subcribed kick channels to sqlite db.

Channel list is inside: constants.go
You can find chatroom id of kick channels from network requests.

How to find chatroom id of channels:
Open DevTools -> Network tab.
Refresh while on channel page
Find https://kick.com/api/v1/channels/{channelname} request.

Inside the json file, chatroom.id is the id of chatroom.

Deployment:

Install docker.
Run below on project directory.

`docker-compose up -d --no-deps --build kick-chat-listener_app`

or install golang and run via:

`go run .`


Db file will be created inside /local-sqlite-data folder.

# Performance

⚙️ Parallel Message Handling

Every incoming message is processed in its own goroutine:
`go UnmarshallAndSendToChannel(...)`
This ensures:
-	Fast processing of high-throughput channels
-	Non-blocking behavior for the main WebSocket reader

📦 Centralized Batching for DB Inserts
A single dedicated goroutine handles batched inserts:
`go HandleDataChannelInserts(db, dataChannel)`
-	Messages are collected via a buffered channel
-	Inserts happen in batches (default size: 10)
-	Reduces DB overhead and improves performance
-	Sqlite compatible architecture
