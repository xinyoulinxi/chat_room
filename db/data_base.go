package chat_db

import (
	"log/slog"
	chat_type "web_server/type"
)

// SaveRoomNameListToFile 将聊天室列表保存到本地文件
func SaveRoomNameListToFile(chatRoomList chat_type.RoomList) {
	err := Default().Set("room_list", &chatRoomList)
	if err != nil {
		slog.Error("Failed to save room list", "error", err)
		return
	}
}

// LoadRoomNameListFromFile 从本地文件中读取聊天室列表
func LoadRoomNameListFromFile() chat_type.RoomList {
	roomList := chat_type.RoomList{}
	_, err := Default().Get("room_list", &roomList)
	if err != nil {
		slog.Error("Failed to read room list", "error", err)
		return roomList
	}
	return roomList
}

// initDefaultChatRoom 初始化默认聊天室
func initDefaultChatRoom(chatName string) *chat_type.ChatRoom {
	chatRoom := chat_type.ChatRoom{RoomName: chatName}
	chatRoom.Messages = make([]chat_type.Message, 0)
	chatRoom.Messages = append(chatRoom.Messages, chat_type.Message{Type: "text", Content: "welcome to " + chatName + "!"})
	return &chatRoom
}

// LoadChatRoomFromLocalFile 从本地文件中读取聊天室信息
func LoadChatRoomFromLocalFile(chatName string) *chat_type.ChatRoom {
	messages := chat_type.Messages{}
	_, err := Default().Get("chatroom_"+chatName, &messages, "chatroom")
	if err != nil {
		slog.Error("Failed to read room message", "room", chatName, "error", err)
		return initDefaultChatRoom(chatName)
	}
	return &chat_type.ChatRoom{RoomName: chatName, Messages: messages}
}

// WriteChatInfoToLocalFile 将聊天室信息保存到本地文件
func WriteChatInfoToLocalFile(chatRoom *chat_type.ChatRoom) error {
	chatName := chatRoom.RoomName
	messages := chatRoom.Messages
	slog.Info("WriteChatInfoToLocalFile", "room", chatName, "messages", messages)
	err := Default().Set("chatroom_"+chatName, &messages, "chatroom")
	if err != nil {
		slog.Error("Failed to write room message", "room", chatName, "error", err)
		return err
	}
	return nil
}

func WriteUsersToLocalFile(users chat_type.Users) error {
	err := Default().Set("user", &users)
	if err != nil {
		slog.Error("Failed to save users", "error", err)
		return err
	}
	return nil
}

func LoadUsersFromLocalFile() chat_type.Users {
	users := chat_type.Users{}
	_, err := Default().Get("user", &users)
	if err != nil {
		slog.Error("Failed to read users", "error", err)
	}
	return users
}
