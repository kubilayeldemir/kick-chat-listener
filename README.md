Listens multiple channel at the same time, saves every message sent to subcribed slack channels to sqlite db.


Channel list is inside: constants.go
You can find chatroom id from network requests.

For example:
Open DevTools -> Network tab.
Refresh while on channel page
Find https://kick.com/api/v1/channels/{channelname} request.

chatroom.id is the id of chatroom.



Deployment:

Install docker.
Run below on project directory.

docker-compose up -d --no-deps --build kick-chat-listener_app


Db file will be created inside /local-sqlite-data folder.
