package models

type Message struct {
	ID           int    `json:"id"`
	Sender       int    `json:"sender"`
	Receiver     int    `json:"receiver"`
	Message_line string `json:"message_line"`
	Shown        bool   `json:"shown"`
}

type Dialog struct {
	Messages []Message `json:"messages"`
	Userhash int       `json:"userhash"`
}

//easyjson:json
type Dialoges []Dialog

//easyjson:json
type Messages []Message
