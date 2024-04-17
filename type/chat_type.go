package chat_type

type ChatRoom struct {
	RoomName string
	Messages []Message
}

type Message struct {
	// MessageType string `json:"type"`
	MsgID    string `json:"id"`
	Type     string `json:"type"` // text  image  roomList file over
	UserName string `json:"userName"`
	Content  string `json:"content"`
	Image    string `json:"image,omitempty"` // Base64-encoded image data
	File     string `json:"file,omitempty"`  // Base64-encoded file data
	SendTime string `json:"sendTime"`
	RoomName string `json:"roomName,omitempty"`

	ChatRoomList []string `json:"chatRoomList,omitempty"`
}

type User struct {
	UserID   string `json:"userId"`
	UserName string `json:"userName"`
	PassWord string `json:"passWord"`
}
