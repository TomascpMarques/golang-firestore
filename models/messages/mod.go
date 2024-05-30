package messages

import "time"

type Message struct {
	Sender     string `form:"sender" json:"sender" firestore:"sender"`
	Content    string `form:"content" json:"content" firestore:"content"`
	SenderType string `form:"sender_type" json:"sender_type" firestore:"sender_type"`
}

type NewMessage struct {
	Message
	Room string `form:"room" json:"room" firestore:"room"`
}

type StoredMessage struct {
	Message
	Sent time.Time `form:"sent" json:"sent"`
}
