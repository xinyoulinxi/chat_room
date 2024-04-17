package chat_room

import (
	"log/slog"
	"sync"
	chat_db "web_server/db"
)

var (
	initRoom    sync.Once
	chatRoomHub = make(map[string]*Room)
	mux         sync.Mutex
)

// InitChatRoomHub 从本地文件中初始化房间信息
func InitChatRoomHub() {
	initRoom.Do(func() {
		roomList := chat_db.LoadRoomNameListFromFile()
		for _, room := range roomList {
			chatRoomHub[room] = getRoom(room)
		}
	})
}

// ListChatRoom 返回当前服务器房间名称列表
func ListChatRoom() []string {
	var chatRoomList = make([]string, 0, len(chatRoomHub))
	for roomName := range chatRoomHub {
		chatRoomList = append(chatRoomList, roomName)
	}
	return chatRoomList
}

func ChatRoomExist(roomName string) bool {
	_, ok := chatRoomHub[roomName]
	return ok
}

// getRoom 从本地文件加载房间信息
func getRoom(roomName string) *Room {
	chatRoom := chat_db.LoadChatRoomFromLocalFile(roomName)
	room := newRoom(chatRoom)
	room.Serve()
	return room
}

// GetChatRoom 获取房间信息,返回房间信息和房间是否之前已存在
func GetChatRoom(roomName string) (*Room, bool) {
	if room, ok := chatRoomHub[roomName]; ok {
		return room, true
	} else {
		mux.Lock()
		defer mux.Unlock()
		if room, ok := chatRoomHub[roomName]; ok {
			return room, true
		}
		room = getRoom(roomName)
		chatRoomHub[roomName] = room
		slog.Info("add new chat room", "roomName", roomName, "chatRoomList", ListChatRoom())
		// 将新房间信息写入本地
		chat_db.SaveRoomNameListToFile(ListChatRoom())
		return room, false
	}
}
