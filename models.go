package main

// Message represents the top-level structure of the JSON input.
type Message struct {
	Event   string `json:"event"`
	Data    string `json:"data"`
	Channel string `json:"channel"`
}

// Data represents the structure of the `data` JSON string.
type Data struct {
	//ID         string `json:"id"` //rarely received number
	ChatroomID int    `json:"chatroom_id"`
	Content    string `json:"content"`
	Type       string `json:"type"`
	CreatedAt  string `json:"created_at"`
	Sender     Sender `json:"sender"`
}

// Sender represents the nested `sender` structure in `Data`.
type Sender struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Slug     string `json:"slug"`
	//Identity Identity `json:"identity"`
}

// Identity represents the nested `identity` structure in `Sender`.
type Identity struct {
	Color  string  `json:"color"`
	Badges []Badge `json:"badges"`
}

// Badge represents the elements in the `badges` array in `Identity`.
type Badge struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
