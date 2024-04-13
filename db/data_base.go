package chat_db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	chat_type "web_server/type"
	utils "web_server/utils"
)

func SaveRoomNameListToFile(chatRoomList []string) {
	jsonData, err := json.Marshal(chatRoomList)
	if err != nil {
		fmt.Println("Failed to marshal room list:", err)
		return
	}
	err = ioutil.WriteFile(utils.RoomListPath, jsonData, 0644)
	if err != nil {
		fmt.Println("Failed to write room list file:", err)
		return
	}
	fmt.Println("Room list saved to file", string(jsonData))
}

func LoadRoomNameListFromFile() []string {
	data, err := ioutil.ReadFile(utils.RoomListPath)
	if err != nil {
		fmt.Println("Failed to read room list file:", err)
		return nil
	}
	var roomList []string
	err = json.Unmarshal(data, &roomList)
	if err != nil {
		fmt.Println("Failed to unmarshal room list:", err)
		return nil
	}
	return roomList
}

func initDefaultChatRoom(chatName string) chat_type.ChatRoom {
	chatRoom := chat_type.ChatRoom{RoomName: chatName}
	chatRoom.Messages = make([]chat_type.Message, 0)
	chatRoom.Messages = append(chatRoom.Messages, chat_type.Message{Type: "text", Content: "welcome to " + chatName + "!"})
	return chatRoom
}

func LoadChatRoomFromLocalFile(chatName string) chat_type.ChatRoom {
	// Load chat room from a file or new a empty chat room
	data, err := ioutil.ReadFile(utils.GetChatRoomFilePath(chatName))
	if err != nil {
		fmt.Println("Failed to read messages file:", err)
		return initDefaultChatRoom(chatName)
	}
	messages := make([]chat_type.Message, 0)
	err = json.Unmarshal(data, &messages)
	if err != nil {
		fmt.Println("Failed to unmarshal messages:", err)
		return initDefaultChatRoom(chatName)
	}
	chatRoom := chat_type.ChatRoom{RoomName: chatName}
	chatRoom.Messages = messages
	return chatRoom
}

func WriteChatInfoToLocalFile(chatRoom *chat_type.ChatRoom) error {
	// Save messages to a file
	jsonData, err := json.Marshal(chatRoom.Messages)
	if err != nil {
		fmt.Println("Failed to marshal messages:", err)
		return err
	}
	err = ioutil.WriteFile(utils.GetChatRoomFilePath(chatRoom.RoomName), jsonData, 0644)
	if err != nil {
		fmt.Println("Failed to write messages file:", err)
		return err
	}
	fmt.Println("Message saved to file")
	return nil
}
