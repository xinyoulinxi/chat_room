package chat_room

import (
	"context"
	"github.com/gorilla/websocket"
	"log/slog"
	"sync"
	chat_db "web_server/db"
	chat_type "web_server/type"
)

const maxHistoryCount = 100

type Room struct {
	ctx  context.Context
	stop context.CancelFunc
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
	*chat_type.ChatRoom
}

func newRoom(room *chat_type.ChatRoom) *Room {
	ctx, cancel := context.WithCancel(context.Background())
	return &Room{
		ctx:        ctx,
		stop:       cancel,
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

func (h *Room) onClientMessage(u *chat_type.User, m chat_type.Message) error {
	h.BroadCast(m)
	return nil
}

func (h *Room) onClientLeave(c *Client) {
	slog.Info("user leave", "id", c.UserID, "userName", c.UserName, "roomName", h.RoomName)
	count := len(h.clients)
	h.unregister <- c
	h.broadRoomUserCountMessage(count - 1)
}

// UserJoin 将用户加入房间
func (h *Room) UserJoin(conn *websocket.Conn, user *chat_type.User) {
	client := newClient(h.ctx, conn, user, h.onClientMessage, h.onClientLeave)
	client.Serve()
	slog.Info("new user join", "id", user.UserID, "userName", user.UserName, "roomName", h.RoomName)
	count := len(h.clients)
	h.register <- client
	h.broadRoomUserCountMessage(count + 1)
}

func (h *Room) broadRoomUserCountMessage(count int) {
	if count <= 0 {
		slog.Info("room user is empty skip broadcast", "roomName", h.RoomName)
		return
	}
	slog.Info("broadcast room user count", "roomName", h.RoomName)
	slog.Info("broadcast room user count end", "roomName", h.RoomName, "userCount", count)
	h.BroadCast(chat_type.NewUserCountMessage(count))
}

func (h *Room) serve() {
	for {
		select {
		case <-h.ctx.Done():
			close(h.register)
			close(h.unregister)
			close(h.broadcast)
			return
		case client := <-h.register:
			h.clients[client] = true
			slog.Info("new user register", "id", client.UserID, "userName", client.UserName, "roomName", h.RoomName, "clientCount", len(h.clients))
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Stop()
			}
			if len(h.clients) == 0 {
				slog.Warn("room is empty", "roomName", h.RoomName)
				RemoveChatRoom(h.RoomName)
			}
		case message := <-h.broadcast:
			if h.Messages.Append(message) {
				_ = chat_db.WriteRoomMessage(h.ChatRoom.RoomName, h.Messages)
			}
			messages := chat_type.Messages{message}
			bytes, err := messages.Serialize()
			if err != nil {
				slog.Error("Failed to convert m to JSON when broadcast message", "roomName", h.RoomName, "error", err)
				continue
			}

			for client := range h.clients {
				_ = client.Send(bytes)
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

func (h *Room) UserCount() int {
	return len(h.clients)
}

func (h *Room) Stop() {
	h.stop()
}
