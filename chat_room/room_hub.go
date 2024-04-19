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
		roomList := chat_db.LoadRoomNameList()
		for _, room := range roomList {
			chatRoomHub[room] = getRoom(room)
			chatRoomList = append(chatRoomList, room)
		}
	})
}

// ListChatRoom 返回当前服务器房间名称列表
func ListChatRoom() []string {
	return chatRoomList
}

func ChatRoomExist(roomName string) (bool, bool) {
	mux.RLock()
	defer mux.RUnlock()
	_, ok := chatRoomHub[roomName]
	if !ok {
		sort.Sort(sort.StringSlice(chatRoomList))
		index := sort.SearchStrings(chatRoomList, roomName)
		if index >= len(chatRoomList) || chatRoomList[index] != roomName {
			return false, false
		} else {
			return true, false
		}
	}
	return true, true
}

// getRoom 从本地文件加载房间信息
func getRoom(roomName string) *Room {
	messages := chat_db.LoadRoomMessage(roomName)
	if messages.Len() == 0 {
		messages.Append(chat_type.Message{Type: "text", Content: "welcome to " + roomName + "!"})
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

		// chatRoomList中查找roomName，如果不存在则添加
		sort.Sort(sort.StringSlice(chatRoomList))
		index := sort.SearchStrings(chatRoomList, roomName)
		if index >= len(chatRoomList) || chatRoomList[index] != roomName {
			chatRoomList = append(chatRoomList, roomName)
		}
		slog.Info("add new chat room", "roomName", roomName, "chatRoomList", ListChatRoom())
		// 将新房间信息写入本地
		chat_db.SaveRoomNameList(ListChatRoom())
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
