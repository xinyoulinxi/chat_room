package chat_db

import (
	"encoding/json"
	"log/slog"
	"os"
	chat_type "web_server/type"
	"web_server/utils"
)

// SaveRoomNameListToFile 将聊天室列表保存到本地文件
func SaveRoomNameListToFile(chatRoomList []string) {
	jsonData, err := json.Marshal(chatRoomList)
	if err != nil {
		slog.Error("Failed to marshal room list", "error", err)
		return
	}
	err = os.WriteFile(utils.RoomListPath, jsonData, 0644)
	if err != nil {
		slog.Error("Failed to write room list file", "error", err)
		return
	}
	slog.Info("ChatRoom list saved to file")
}

// LoadRoomNameListFromFile 从本地文件中读取聊天室列表
func LoadRoomNameListFromFile() []string {
	data, err := os.ReadFile(utils.RoomListPath)
	var roomList []string
	if err != nil {
		slog.Error("Failed to read room list file", "error", err)
		return roomList
	}
	err = json.Unmarshal(data, &roomList)
	if err != nil {
		slog.Error("Failed to unmarshal room list", "error", err)
		return roomList
	}
	return roomList
}

// initDefaultChatRoom 初始化默认聊天室
func initDefaultChatRoom(chatName string) chat_type.ChatRoom {
	chatRoom := chat_type.ChatRoom{RoomName: chatName}
	chatRoom.Messages = make([]chat_type.Message, 0)
	chatRoom.Messages = append(chatRoom.Messages, chat_type.Message{Type: "text", Content: "welcome to " + chatName + "!"})
	return chatRoom
}

// LoadChatRoomFromLocalFile 从本地文件中读取聊天室信息
func LoadChatRoomFromLocalFile(chatName string) chat_type.ChatRoom {
	// Load chat room from a file or new a empty chat room
	data, err := os.ReadFile(utils.GetChatRoomFilePath(chatName))
	if err != nil {
		slog.Error("Failed to read messages file", "error", err)
		return initDefaultChatRoom(chatName)
	}
	messages := make([]chat_type.Message, 0)
	err = json.Unmarshal(data, &messages)
	if err != nil {
		slog.Error("Failed to unmarshal messages", "error", err)
		return initDefaultChatRoom(chatName)
	}
	return chat_type.ChatRoom{RoomName: chatName, Messages: messages}
}

// WriteChatInfoToLocalFile 将聊天室信息保存到本地文件
func WriteChatInfoToLocalFile(chatRoom *chat_type.ChatRoom) error {
	// Save messages to a file
	jsonData, err := json.Marshal(chatRoom.Messages)
	if err != nil {
		slog.Error("Failed to marshal messages", "error", err)
		return err
	}
	err = os.WriteFile(utils.GetChatRoomFilePath(chatRoom.RoomName), jsonData, 0644)
	if err != nil {
		slog.Error("Failed to write messages file", "error", err)
		return err
	}
	slog.Info("Message saved to file")
	return nil
}

func WriteUsersToLocalFile(user []chat_type.User) error {
	// Save messages to a file
	jsonData, err := json.Marshal(user)
	if err != nil {
		slog.Error("Failed to marshal user", "error", err)
		return err
	}
	err = os.WriteFile(utils.UserListPath, jsonData, 0644)
	if err != nil {
		slog.Error("Failed to write user file", "error", err)
		return err
	}
	slog.Info("User saved to file")
	return nil
}

func LoadUsersFromLocalFile() []chat_type.User {
	data, err := os.ReadFile(utils.UserListPath)
	var users []chat_type.User
	if err != nil {
		slog.Error("Failed to read user file", "error", err)
		return users
	}
	err = json.Unmarshal(data, &users)
	if err != nil {
		slog.Error("Failed to unmarshal user", "error", err)
		return users
	}
	return users
}
