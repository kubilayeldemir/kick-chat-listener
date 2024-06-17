Listens multiple channels at the same time, saves every message sent to subcribed slack channels to sqlite db.

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
