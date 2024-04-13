package chat_room

import (
	"encoding/json"
	"fmt"
	"net/http"
	chat_db "web_server/db"
)

func GetChatRoomList() []string {
	chatRoomList := chat_db.LoadRoomNameListFromFile()
	return chatRoomList
}

func ChatRoomListHandler(w http.ResponseWriter, r *http.Request) {
	chatRoomList = chat_db.LoadRoomNameListFromFile()
	fmt.Println("ChatRoomListHandler:", chatRoomList)
	// Convert chat room list to JSON
	jsonMsg, err := json.Marshal(chatRoomList)
	if err != nil {
		fmt.Println("Failed to convert message to JSON:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonMsg)
}
