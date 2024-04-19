package chat_room

import (
	"log/slog"
	"sync"
	chat_db "web_server/db"
	chat_type "web_server/type"
)

var (
	initRoom    sync.Once
	chatRoomHub = make(map[string]*Room)
	mux         sync.RWMutex
)

// InitChatRoomHub 从本地文件中初始化房间信息
func InitChatRoomHub() {
	initRoom.Do(func() {
		roomList := chat_db.LoadRoomNameList()
		for _, room := range roomList {
			chatRoomHub[room] = getRoom(room)
		}
	})
}

// getRoom 从本地文件加载房间信息
func getRoom(roomName string) *Room {
	messages := chat_db.LoadRoomMessage(roomName)
	if messages.Len() == 0 {
		messages.Append(chat_type.NewNoticeMessage("welcome to " + roomName + "!"))
	}
	room := newRoom(&chat_type.ChatRoom{RoomName: roomName, Messages: messages})
	room.Serve()
	return room
}

// GetChatRoom 获取房间信息,返回房间信息和房间是否之前已存在
func GetChatRoom(roomName string) (*Room, bool) {
	mux.RLock()
	if room, ok := chatRoomHub[roomName]; ok {
		mux.RUnlock()
		return room, true
	} else {
		mux.RUnlock()
		mux.Lock()
		defer mux.Unlock()
		if room, ok := chatRoomHub[roomName]; ok {
			return room, true
		}
		room = getRoom(roomName)
		chatRoomHub[roomName] = room

		// 添加room到列表中
		if ok, err := chat_db.AppendRoomName(roomName); err != nil {
			slog.Error("failed to append room name", "roomName", roomName, "error", err)
		} else if ok {
			slog.Info("add new chat room", "roomName", roomName, "chatRoomList", chat_db.LoadRoomNameList())
		}
		roomList := chat_db.LoadRoomNameList()
		for name, room := range chatRoomHub {
			if name == roomName {
				continue
			}
			room.BroadCast(chat_type.NewRoomListMessage(roomList))
		}
		return room, false
	}
}

func RemoveChatRoom(roomName string) {
	mux.Lock()
	defer mux.Unlock()
	if room, ok := chatRoomHub[roomName]; ok {
		delete(chatRoomHub, roomName)
		slog.Info("remove chat room", "roomName", roomName)
		room.Stop()
	}
}
