package chat_type

import "encoding/json"

type Message struct {
	MsgID        string          `json:"id"`                     // msg id for every persistent message
	SendTime     string          `json:"sendTime"`               // 'YYYY-MM-DD HH:mm:ss' format for time
	Type         string          `json:"type"`                   // text image roomList file over userCount
	Content      string          `json:"content"`                // filename for image or file, and content for text
	Image        string          `json:"image,omitempty"`        // src for type image
	File         string          `json:"file,omitempty"`         // src for type file
	RoomName     string          `json:"roomName,omitempty"`     // room name for text, image or file
	UserName     string          `json:"userName"`               // username for text, image or file
	AvatarUrl    string          `json:"avatarUrl,omitempty"`    // avatar url for text, image or file
	Data         json.RawMessage `json:"data,omitempty"`         // for type userCount
	ChatRoomList []string        `json:"chatRoomList,omitempty"` // for type roomList
}

type Messages []Message

func (m *Messages) Len() int {
	return len(*m)
}

func (m *Messages) Append(msg Message) {
	*m = append(*m, msg)
}

func (m *Messages) LastN(n int) []Message {
	if len(*m) > n {
		return (*m)[len(*m)-n:]
	} else {
		return *m
	}
}
