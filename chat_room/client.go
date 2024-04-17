package chat_room

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log/slog"
	"sync"
	"time"
	chat_type "web_server/type"
	"web_server/utils"
)

const (
	pingPeriod     = 60 * time.Second
	imageSizeLimit = 1024 * 1024 * 20
)

type onMessage func(user *chat_type.User, message chat_type.Message) error
type onClientLeave func(client *Client)

type Client struct {
	exit bool
	init sync.Once
	conn *websocket.Conn
	send chan []byte
	*chat_type.User

	// 客户端消息到达
	onMessage onMessage
	// 客户端退出
	onClientLeave onClientLeave
}

func (c *Client) Serve() {
	c.init.Do(func() {
		slog.Info("Client Serve", "id", c.UserID, "userName", c.UserName)
		go c.readPump()
		go c.writePump()
	})
}

func (c *Client) Send(m chat_type.Message) error {
	slog.Info("sendMessage to client", "id", c.UserID, "userName", c.UserName, "content", m.Content, "type", m.Type, "roomName", m.RoomName, "sendTime", m.SendTime)
	jsonMsg, err := json.Marshal([]chat_type.Message{m})
	if err != nil {
		slog.Error("Failed to convert m to JSON", "error", err)
		return err
	}
	c.send <- jsonMsg
	return nil
}

func (c *Client) readPump() {
	defer func() {
		c.exit = true
		c.onClientLeave(c)
	}()
	for {
		if c.exit {
			slog.Warn("Client exit, stop read", "id", c.UserID, "userName", c.UserName)
			break
		}
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("Client unexpected close", "id", c.UserID, "userName", c.UserName, "error", err)
			} else {
				slog.Error("Failed to read message from WebSocket", "id", c.UserID, "userName", c.UserName, "error", err)
			}
			break
		}

		var message chat_type.Message
		err = json.Unmarshal(msg, &message)
		if err != nil {
			slog.Error("Failed to parse message", "error", err)
			continue
		}
		// 限制message.Image的文件大小
		if message.Image != "" {
			if len(message.Image) > imageSizeLimit {
				message.Image = ""
			}
		}
		utils.TryTransferImagePathToMessage(&message)
		// get tUer
		// tUer := user.GetUserById(id)
		// 暂时直接使用消息带上来的userName
		// message.UserName = tUer.UserName
		message.MsgID = utils.GenerateId()
		message.SendTime = utils.GetCurTime()
		// 广播消息
		err = c.onMessage(c.User, message)
	}
}

func (c *Client) writePump() {
	slog.Info("Client writePump", "id", c.UserID, "userName", c.UserName)
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()
	for {
		if c.exit {
			slog.Warn("Client exit, stop write", "id", c.UserID, "userName", c.UserName)
			break
		}
		select {
		case message, ok := <-c.send:
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				slog.Warn("The hub room closed the channel")
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				slog.Error("Failed to send message", "id", c.UserID, "userName", c.UserName, "message", string(message), "error", err)
			}
		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				slog.Error("Failed to send ping message", "id", c.UserID, "userName", c.UserName, "error", err)
				return
			}
		}
	}
	slog.Info("Client writePump exit", "id", c.UserID, "userName", c.UserName)
}
