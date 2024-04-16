package chat_room

import (
	"encoding/json"
	"log/slog"
	"net/http"
	chat_db "web_server/db"
)

func GetChatRoomList() []string {
	chatRoomList := chat_db.LoadRoomNameListFromFile()
	return chatRoomList
}

func ChatRoomListHandler(w http.ResponseWriter, r *http.Request) {
	chatRoomList = chat_db.LoadRoomNameListFromFile()
	slog.Info("ChatRoomListHandler", "chatRoomList", chatRoomList)
	// Convert chat room list to JSON
	jsonMsg, err := json.Marshal(chatRoomList)
	if err != nil {
		slog.Error("Failed to convert message to JSON", "error", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonMsg)
}
