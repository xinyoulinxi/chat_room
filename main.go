package main

import (
	"log/slog"
	"os"
	"web_server/chat_room"
	chat_db "web_server/db"
	"web_server/page_handler"
	"web_server/user"
	"web_server/utils"
)

func main() {
	utils.InitEnv()

	if storage, err := chat_db.NewFileStorage(utils.DbPath); err != nil {
		slog.Error("init file storage error", "error", err)
		os.Exit(1)
	} else {
		chat_db.Init(storage)
	}

	chat_room.InitChatRoomHub()
	user.InitUserInfos()
	page_handler.StartWebServer()
}
