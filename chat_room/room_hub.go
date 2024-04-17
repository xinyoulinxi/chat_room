package chat_room

import (
	"log/slog"
	"sort"
	"sync"
	chat_db "web_server/db"
	chat_type "web_server/type"
)

var (
	initRoom     sync.Once
	chatRoomHub  = make(map[string]*Room)
	chatRoomList = make([]string, 0)
	mux          sync.RWMutex
)

// InitChatRoomHub 从本地文件中初始化房间信息
func InitChatRoomHub() {
	initRoom.Do(func() {
		roomList := chat_db.LoadRoomNameListFromFile()
		for _, room := range roomList {
			chatRoomHub[room] = getRoom(room)
			chatRoomList = append(chatRoomList, room)
		}
	})
}

// ListChatRoom 返回当前服务器房间名称列表
func ListChatRoom() []string {
	// mux.RLock()
	// defer mux.RUnlock()
	// var chatRoomList = make([]string, 0, len(chatRoomHub))
	// for roomName := range chatRoomHub {
	// 	chatRoomList = append(chatRoomList, roomName)
	// }
	// return chatRoomList
	return chatRoomList
}

func ChatRoomExist(roomName string) bool {
	mux.RLock()
	defer mux.RUnlock()
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

		// chatRoomList中查找roomName，如果不存在则添加
		sort.Sort(sort.StringSlice(chatRoomList))
		index := sort.SearchStrings(chatRoomList, roomName)
		if index >= len(chatRoomList) || chatRoomList[index] != roomName {
			chatRoomList = append(chatRoomList, roomName)
		}
		slog.Info("add new chat room", "roomName", roomName, "chatRoomList", ListChatRoom())
		// 将新房间信息写入本地
		chat_db.SaveRoomNameListToFile(ListChatRoom())
		for name, room := range chatRoomHub {
			if name == roomName {
				continue
			}
			room.BroadCast(chat_type.Message{Type: "roomList", ChatRoomList: ListChatRoom()})
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
