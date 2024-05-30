package rooms

type Room struct {
	Name        string `form:"name" json:"name" firestore:"name"`
	Owner       string `form:"owner" json:"owner" firestore:"owner"`
	Description string `form:"description" json:"description" firestore:"description"`
	Category    string `form:"category" json:"category" firestore:"category"`
}

type RoomDefined struct {
	Room
	MsgCount int `form:"msg_count" json:"msg_count" firestore:"msg_count"`
}
