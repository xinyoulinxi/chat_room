package chat_type

import (
	"encoding/json"
)

type RoomList []string

func (r *RoomList) Serialize() ([]byte, error) {
	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func (r *RoomList) Deserialize(bytes []byte) error {
	return json.Unmarshal(bytes, &r)
}

type Messages []Message

func (m *Messages) Serialize() ([]byte, error) {
	jsonData, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}
func (m *Messages) Deserialize(bytes []byte) error {
	return json.Unmarshal(bytes, &m)
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

type ChatRoom struct {
	RoomName string
	Messages Messages
}

type User struct {
	UserID   string `json:"userId"`
	UserName string `json:"userName"`
	PassWord string `json:"passWord"`
}

type Users []User

func (u *Users) Serialize() ([]byte, error) {
	jsonData, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}
func (u *Users) Deserialize(bytes []byte) error {
	return json.Unmarshal(bytes, &u)
}
