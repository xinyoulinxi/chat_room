package main

import (
	page_handler "web_server/page_handler"
	utils "web_server/utils"
)

func main() {
	utils.InitEnv()
	page_handler.StartWebServer()
}
