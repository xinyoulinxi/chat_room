package chat_room

import (
	"github.com/gorilla/websocket"
	"log/slog"
	"sync"
	chat_type "web_server/type"
)

const maxHistoryCount = 100

type Room struct {
	init sync.Once
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan chat_type.Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Current room
	chat_type.ChatRoom
}

func newRoom(room chat_type.ChatRoom) *Room {
	return &Room{
		broadcast:  make(chan chat_type.Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		ChatRoom:   room,
	}
}

// BroadCast 全房间广播消息
func (h *Room) BroadCast(m chat_type.Message) {
	h.broadcast <- m
}

func (h *Room) UserJoin(conn *websocket.Conn, user *chat_type.User) {
	client := &Client{
		User: user,
		conn: conn,
		onMessage: func(u *chat_type.User, m chat_type.Message) error {
			h.BroadCast(m)
			return nil
		},
		send: make(chan []byte),
		onClientLeave: func(c *Client) {
			h.unregister <- c
		},
	}
	client.Serve()
	slog.Info("new user join", "id", user.UserID, "userName", user.UserName, "roomName", h.RoomName, "memberSize", len(h.clients))
	h.register <- client
}

// sendHistory 发送历史消息
func (h *Room) sendHistory(c *Client) {
	var messages []chat_type.Message
	if len(h.Messages) > maxHistoryCount {
		// 保留最新100条
		messages = h.Messages[len(h.Messages)-maxHistoryCount:]
	} else {
		messages = h.Messages
	}

	for _, message := range messages {
		_ = c.Send(message)
	}
	_ = c.Send(chat_type.Message{Type: "over", RoomName: h.RoomName})
}

func (h *Room) serve() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.sendHistory(client)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			h.Messages = append(h.Messages, message)
			// TODO 消息持久化
			for client := range h.clients {
				_ = client.Send(message)
			}
		}
	}
}

func (h *Room) Serve() {
	h.init.Do(func() {
		slog.Info("room serve", "roomName", h.RoomName)
		go h.serve()
	})
}
