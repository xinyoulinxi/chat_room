package chat_type

import "github.com/gorilla/websocket"

type ChatRoom struct {
	RoomName    string
	Connections []*websocket.Conn
	Messages    []Message
	Users       []string
}

type Message struct {
	// MessageType string `json:"type"`
	Type     string `json:"type"` // text or image or roomList
	UserID   string `json:"userId"`
	Content  string `json:"content"`
	Image    string `json:"image,omitempty"` // Base64-encoded image data
	SendTime string `json:"sendTime"`
	RoomName string `json:"roomName,omitempty"`

	ChatRoomList []string `json:"chatRoomList,omitempty"`
}
