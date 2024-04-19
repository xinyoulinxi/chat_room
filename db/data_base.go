package chat_db

import (
	"log/slog"
	chat_type "web_server/type"
)

var (
	roomStorage    RoomStorage
	messageStorage MessageStorage
	userStorage    UserStorage
)

func Init(h Handler) {
	roomStorage = RoomStorage{handler: h}
	messageStorage = MessageStorage{handler: h}
	userStorage = UserStorage{handler: h}
}

// SaveRoomNameList 将聊天室列表保存到本地文件
func SaveRoomNameList(chatRoomList chat_type.RoomList) {
	err := roomStorage.SaveAll(chatRoomList)
	if err != nil {
		slog.Error("Failed to save room list", "error", err)
		return
	}
}

// LoadRoomNameList 从本地文件中读取聊天室列表
func LoadRoomNameList() chat_type.RoomList {
	list, err := roomStorage.LoadAll()
	if err != nil {
		slog.Error("Failed to read room list", "error", err)
		return list
	}
	return list
}

func CheckRoomName(chatName string) bool {
	return roomStorage.FindRoom(chatName)
}

func AppendRoomName(chatName string) (bool, error) {
	return roomStorage.Append(chatName)
}

// LoadRoomMessage 从本地文件中读取聊天室消息历史
func LoadRoomMessage(chatName string) chat_type.Messages {
	messages, err := messageStorage.LoadAll(chatName)
	if err != nil {
		slog.Error("Failed to read room message", "room", chatName, "error", err)
	}
	return messages
}

// WriteRoomMessage 将聊天室消息历史保存到本地文件
func WriteRoomMessage(chatName string, messages chat_type.Messages) error {
	slog.Info("WriteRoomMessage", "room", chatName, "messages", messages)
	err := messageStorage.SaveAll(chatName, messages)
	if err != nil {
		slog.Error("Failed to write room message", "room", chatName, "error", err)
		return err
	}
	return nil
}

func WriteUsers(users chat_type.Users) error {
	err := userStorage.SaveAll(users)
	if err != nil {
		slog.Error("Failed to save users", "error", err)
		return err
	}
	return nil
}

func LoadUsers() chat_type.Users {
	users, err := userStorage.LoadAll()
	if err != nil {
		slog.Error("Failed to read users", "error", err)
	}
	return users
}
