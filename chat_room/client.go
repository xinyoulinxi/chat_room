package chat_room

import (
	"context"
	"github.com/gorilla/websocket"
	"log/slog"
	"sync"
	"time"
	chat_type "web_server/type"
)

const (
	pingPeriod = 60 * time.Second
)

type onMessage func(user *chat_type.User, message chat_type.Message) error
type onClientLeave func(client *Client)

type Client struct {
	ctx  context.Context
	stop context.CancelFunc
	init sync.Once
	conn *websocket.Conn
	send chan []byte
	*chat_type.User

	// 客户端消息到达
	onMessage onMessage
	// 客户端退出
	onClientLeave onClientLeave
}

func newClient(ctx context.Context, conn *websocket.Conn, user *chat_type.User, onMessage onMessage, onClientLeave onClientLeave) *Client {
	ctx, cancel := context.WithCancel(ctx)
	client := &Client{
		ctx:           ctx,
		stop:          cancel,
		User:          user,
		conn:          conn,
		send:          make(chan []byte),
		onMessage:     onMessage,
		onClientLeave: onClientLeave,
	}
	return client
}

func (c *Client) Serve() {
	c.init.Do(func() {
		slog.Info("Client Serve", "id", c.UserID, "userName", c.UserName)
		go c.readPump()
		go c.writePump()
	})
}

func (c *Client) Stop() {
	c.stop()
}

func (c *Client) Send(m []byte) error {
	if c.send == nil {
		return nil
	} else {
		c.send <- m
	}
	return nil
}

func (c *Client) readPump() {
	defer func() {
		slog.Warn("Client exit, stop read", "id", c.UserID, "userName", c.UserName)
		c.stop()
		c.onClientLeave(c)
	}()
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("Client unexpected close", "id", c.UserID, "userName", c.UserName, "error", err)
			} else {
				slog.Error("Failed to read message from WebSocket", "id", c.UserID, "userName", c.UserName, "error", err)
			}
			return
		}

		var message chat_type.Message
		err = message.Deserialize(msg)
		if err != nil {
			slog.Error("Failed to parse message", "error", err)
			continue
		}
		// get tUer
		// tUer := user.GetUserById(id)
		// 暂时直接使用消息带上来的userName
		// message.UserName = tUer.UserName
		message.UserID = c.UserID
		message.UserName = c.UserName
		message.Wrap()
		// 广播消息
		err = c.onMessage(c.User, message)
	}
}

func (c *Client) writePump() {
	slog.Info("Client writePump", "id", c.UserID, "userName", c.UserName)
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		slog.Info("Client writePump exit", "id", c.UserID, "userName", c.UserName)
		ticker.Stop()
		_ = c.conn.Close()
	}()
	for {
		select {
		case <-c.ctx.Done():
			close(c.send)
			c.send = nil
			return
		case message, ok := <-c.send:
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				slog.Warn("The hub room closed the channel")
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					slog.Error("Client unexpected close", "id", c.UserID, "userName", c.UserName, "error", err)
				} else {
					slog.Error("Failed to send message", "id", c.UserID, "userName", c.UserName, "message", string(message), "error", err)
				}
				return
			}
		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				slog.Error("Failed to send ping message", "id", c.UserID, "userName", c.UserName, "error", err)
				return
			}
		}
	}
}
