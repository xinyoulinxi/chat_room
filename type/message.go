package chat_type

import (
	"encoding/json"
	"errors"
	"github.com/bwmarrin/snowflake"
	"log/slog"
	"math/rand"
	"time"
)

var snowId *snowflake.Node

func init() {
	id, err := snowflake.NewNode(rand.New(rand.NewSource(time.Now().UnixMilli())).Int63n(1024))
	if err != nil {
		panic(err)
	}
	snowId = id
}

type MsgType string

const (
	TextMsg      MsgType = "text"      // 文字消息
	ImageMsg     MsgType = "image"     // 图片消息
	FileMsg      MsgType = "file"      // 文件消息
	NoticeMsg    MsgType = "notice"    // 房间通知消息
	RoomListMsg  MsgType = "roomList"  // 房间列表变动消息
	UserCountMsg MsgType = "userCount" // 用户数量变动消息
	UserListMsg  MsgType = "userList"  // 用户列表变动消息
)

type (
	Message struct {
		MsgID    string  `json:"id"`                 // msg id for every persistent message
		SendTime string  `json:"sendTime"`           // 'YYYY-MM-DD HH:mm:ss' format for time
		Type     MsgType `json:"type"`               // text image roomList file userCount
		Content  string  `json:"content,omitempty"`  // filename for image or file, and content for text
		UserID   string  `json:"userId,omitempty"`   // userid for text, image or file
		UserName string  `json:"userName,omitempty"` // username for text, image or file
		// AvatarUrl string `json:"avatarUrl,omitempty"` // avatar url for text, image or file
		// RoomName string `json:"roomName,omitempty"` // room name for text, image or file

		ImagePart
		FilePart
		ExtPart
	}

	// ImagePart 用于标示图片消息部分字段
	ImagePart struct {
		Image string `json:"image,omitempty"` // src for type image
	}

	// FilePart 用于标示文件消息部分字段
	FilePart struct {
		File string `json:"file,omitempty"` // src for type file
	}

	ExtPart struct {
		Data any `json:"data,omitempty"` // for type userCount
	}
)

func NewUserCountMessage(count int) Message {
	m := Message{
		Type: UserCountMsg,
		ExtPart: ExtPart{
			Data: count,
		},
	}
	m.Wrap()
	return m
}

func NewRoomListMessage(roomList []string) Message {
	m := Message{
		Type: RoomListMsg,
		ExtPart: ExtPart{
			Data: roomList,
		},
	}
	m.Wrap()
	return m
}

func NewUserListMessage(userList []string) Message {
	m := Message{
		Type: UserListMsg,
		ExtPart: ExtPart{
			Data: userList,
		},
	}
	m.Wrap()
	return m
}

func NewNoticeMessage(content string) Message {
	m := Message{
		Type:    NoticeMsg,
		Content: content,
	}
	m.Wrap()
	return m
}

func (m *Message) Wrap() {
	m.MsgID = snowId.Generate().String()
	m.SendTime = time.Now().Format(time.DateTime)
}

func (m *Message) Deserialize(bytes []byte) error {
	err := json.Unmarshal(bytes, &m)
	if err != nil {
		return err
	}
	if m.Type == "" {
		return errors.New("unsupported message type")
	}
	return nil
}

type Messages []Message

func (m *Messages) Len() int {
	return len(*m)
}

func (m *Messages) Append(msg Message) bool {
	switch msg.Type {
	case TextMsg, ImageMsg, FileMsg, NoticeMsg:
		*m = append(*m, msg)
		return true
	}
	return false
}

func (m *Messages) LastN(n int) []Message {
	if len(*m) > n {
		return (*m)[len(*m)-n:]
	} else {
		return *m
	}
}

// 输入index和n，index代表当前需要数组的倒数第index个长度为n的slice，返回这个slice。如果剩下的数量不足n，则返回剩下的
func (m *Messages) IndexN(index int, n int) []Message {
	length := len(*m)
	if length < n && index == 0 {
		return *m
	}
	start := length - (index+1)*n
	end := start + n
	slog.Info("IndexN", "start", start, "len", len(*m))
	if start < 0 {
		if end > 0 {
			return (*m)[0:end]
		}
		return nil
	}
	if end > length {
		return (*m)[start:]
	}
	return (*m)[start:end]
}

func (m *Messages) Serialize() ([]byte, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
