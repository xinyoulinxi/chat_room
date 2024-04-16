package main

import (
	"web_server/page_handler"
	"web_server/utils"
)

func main() {
	utils.InitEnv()
	page_handler.StartWebServer()
}
