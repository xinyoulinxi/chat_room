package main

import (
	"web_server/chat_room"
	"web_server/page_handler"
	"web_server/user"
	"web_server/utils"
)

func main() {
	utils.InitEnv()
	chat_room.InitChatRoomHub()
	user.InitUserInfos()
	page_handler.StartWebServer()
}
