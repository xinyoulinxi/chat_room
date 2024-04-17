package chat_room

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
	chat_type "web_server/type"
	"web_server/user"
	"web_server/utils"
	// "github.com/gorilla/handlers"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func CreateChatRoomHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("id")
	roomName := r.URL.Query().Get("roomName")
	slog.Info("CreateChatRoomHandler", "roomName", roomName, "userId", userId)
	if userId == "" {
		utils.WriteResponse(w, 1, "Invalid user id")
		return
	}
	if user.UserRegisted(userId) == false {
		utils.WriteResponse(w, 1, "User does not exist")
		return
	}
	if roomName == "" {
		utils.WriteResponse(w, 1, "Invalid chat room name")
		return
	}
	_, exist := GetChatRoom(roomName)
	if exist {
		utils.WriteResponse(w, 1, "Chat room already exist")
		return
	}
	utils.WriteResponse(w, 0, "Success")
}

func HistoryMessagesHandler(w http.ResponseWriter, r *http.Request) {
	// 确保关闭请求体
	defer r.Body.Close()
	// 检查请求方法是否为POST
	if r.Method != http.MethodPost {
		utils.WriteResponse(w, chat_type.ErrorInvalidInput, "Invalid request method")
		return
	}

	// 解析请求体
	var loginData struct {
		UserId   string `json:"userId"`
		ChatRoom string `json:"chatRoom"`
		Count    int    `json:"count"`
	}

	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		utils.WriteResponse(w, chat_type.ErrorCodeFail, "Failed to parse request body")
		return
	}
	userId := loginData.UserId
	roomName := loginData.ChatRoom

	slog.Info("HistoryMessagesHandler", "userId", userId, "ChatRoom", roomName)
	if userId == "" || roomName == "" {
		utils.WriteResponse(w, chat_type.ErrorInvalidInput, "Invalid user id or chat room name")
		return
	}

	if !user.UserRegisted(userId) {
		utils.WriteResponse(w, chat_type.ErrorUserNotExist, "User not exist")
		return
	}
	if exist, _ := ChatRoomExist(roomName); !exist {
		utils.WriteResponse(w, chat_type.ErrorInvalidInput, "Chat room not exist")
		return
	}
	chatRoom, _ := GetChatRoom(roomName)
	var messages []chat_type.Message
	count := loginData.Count
	if count <= 0 {
		count = maxHistoryCount
	}
	if len(chatRoom.Messages) > count {
		// 保留最新100条
		messages = chatRoom.Messages[len(chatRoom.Messages)-count:]
	} else {
		messages = chatRoom.Messages
	}
	// 将chatRoom.Messages转换成json字符串
	jsonMsg, err := json.Marshal(messages)
	if err != nil {
		utils.WriteResponse(w, chat_type.ErrorCodeFail, "fail to fetch messages")
		slog.Error("Failed to convert message to JSON", "error", err)
		return
	}
	slog.Info("HistoryMessagesHandler", "messages size", len(messages))
	utils.WriteResponseWithData(w, chat_type.ErrorCodeSuccess, "Success", jsonMsg)
}

func UploadFilehandler(w http.ResponseWriter, r *http.Request) {
	// 确保关闭请求体
	defer r.Body.Close()
	// 检查请求方法是否为POST
	if r.Method != http.MethodPost {
		utils.WriteResponse(w, chat_type.ErrorInvalidInput, "Invalid request method")
		return
	}

	// 解析请求体
	var loginData struct {
		UserId   string `json:"userId"`
		FileName string `json:"filename"`
		Data     string `json:"data,omitempty"` // Base64-encoded data
	}

	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		utils.WriteResponse(w, chat_type.ErrorCodeFail, "Failed to parse request body")
		return
	}
	userId := loginData.UserId
	data := loginData.Data
	fileName := loginData.FileName

	if data != "" {
		if len(data) > 1024*1024*20 {
			utils.WriteResponse(w, chat_type.ErrorInvalidInput, "File too large")
			return
		}
	}
	filePath, fileType, err := utils.SaveFile(fileName, &data)
	if err != nil {
		utils.WriteResponse(w, chat_type.ErrorInvalidInput, "Failed to save file")
		return
	}
	var fileStruct struct {
		FilaPath string `json:"filePath"`
		FileType string `json:"fileType"`
	}
	fileStruct.FilaPath = filePath
	fileStruct.FileType = fileType
	// file path 转换为json
	filePathJson, err := json.Marshal(fileStruct)
	if err != nil {
		utils.WriteResponse(w, chat_type.ErrorInvalidInput, "Failed to marshal file path")
		return
	}
	utils.WriteResponseWithData(w, 0, "Success", filePathJson)
	slog.Info("uploadFilehandler", "userId", userId)
}

func ChatRoomListHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("id")
	slog.Info("ChatRoomListHandler", "userId", userId)
	if userId == "" {
		utils.WriteResponse(w, chat_type.ErrorInvalidInput, "Invalid user id")
		return
	}
	if user.UserRegisted(userId) == false {
		utils.WriteResponse(w, chat_type.ErrorUserNotExist, "User does not exist")
		return
	}

	roomList := ListChatRoom()
	slog.Info("ChatRoomListHandler", "chatRoomList", roomList)
	// Convert chat room list to JSON
	jsonMsg, err := json.Marshal(roomList)
	if err != nil {
		utils.WriteResponse(w, chat_type.ErrorCodeFail, "Failed to get room list")
		slog.Error("Failed to convert message to JSON", "error", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	utils.WriteResponseWithData(w, chat_type.ErrorCodeSuccess, "Success", jsonMsg)
}

// ChatRoomHandler 处理用户加入房间
func ChatRoomHandler(w http.ResponseWriter, r *http.Request) {
	// 获取request中的参数,比如id和chatroom
	id := r.URL.Query().Get("id")
	chatRoomName := r.URL.Query().Get("chatroom")
	if chatRoomName == "" || chatRoomName == "null" {
		chatRoomName = "default"
	}

	// 根据id获取用户
	u := user.GetUserById(id)
	if u == nil {
		utils.WriteResponse(w, chat_type.ErrorInvalidInput, "Invalid user id")
		return
	}
	if exist, _ := ChatRoomExist(chatRoomName); !exist {
		utils.WriteResponse(w, chat_type.ErrorInvalidInput, "Chat room not exist")
		return
	}
	// Get chat room by name
	room, _ := GetChatRoom(chatRoomName)

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade connection to WebSocket", http.StatusInternalServerError)
		return
	}

	room.UserJoin(conn, user.GetUserById(id))
}
