package chat_room

import (
	"encoding/json"
	"log/slog"
	"net/http"
	chat_db "web_server/db"
	chat_type "web_server/type"
	"web_server/user"
	"web_server/utils"

	"github.com/gorilla/websocket"
	// "github.com/gorilla/handlers"
)

const maxHistoryCount = 100

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var chatRoomMap = make(map[string]*chat_type.ChatRoom)
var chatRoomList = make([]string, 0)

func getChatRoomByName(roomName string) *chat_type.ChatRoom {
	if chatRoom, ok := chatRoomMap[roomName]; ok {
		return chatRoom
	} else {
		chatRoom := chat_db.LoadChatRoomFromLocalFile(roomName)
		chatRoomMap[roomName] = &chatRoom
		isExist := false
		for _, room := range chatRoomList {
			if room == roomName {
				isExist = true
				break
			}
		}
		if !isExist {
			slog.Info("add new chat room", "roomName", roomName, "chatRoomList", chatRoomList)
			chatRoomList = append(chatRoomList, roomName)
		}
		return &chatRoom
	}
}

func printMessageInfo(message *chat_type.Message) {
	slog.Info("printMessageInfo", "id", message.MsgID, "type", message.Type, "userName", message.UserName, "content", message.Content, "image", message.Image, "file", message.File, "sendTime", message.SendTime, "roomName", message.RoomName)
}

func CloseChatRoom(chatRoom *chat_type.ChatRoom) {
	chat_db.WriteChatInfoToLocalFile(chatRoom)
	delete(chatRoomMap, chatRoom.RoomName)
}

func removeConnection(conn *websocket.Conn, connections []*websocket.Conn) []*websocket.Conn {
	for i, c := range connections {
		if c == conn {
			return append(connections[:i], connections[i+1:]...)
		}
	}
	return connections
}
func addNewUserToChatRoom(chatRoom *chat_type.ChatRoom, id string) {
	for _, user := range chatRoom.Users {
		if user == id {
			return
		}
	}
	chatRoom.Users = append(chatRoom.Users, id)
}

func sendMessage(message chat_type.Message, c *websocket.Conn) error {
	// Convert Message struct to JSON
	slog.Info("sendMessage", "username", message.UserName, "content", message.Content, "type", message.Type, "roomName", message.RoomName, "sendTime", message.SendTime)
	jsonMsg, err := json.Marshal([]chat_type.Message{message})
	if err != nil {
		slog.Error("Failed to convert message to JSON", "error", err)
		return err
	}

	// Write message to WebSocket
	err = c.WriteMessage(websocket.TextMessage, jsonMsg)
	return err
}

func CreateChatRoomHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request parameters
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
	for _, room := range chatRoomList {
		if room == roomName {
			utils.WriteResponse(w, 1, "Chat room already exist")
			return
		}
	}

	chatRoom := getChatRoomByName(roomName)
	chat_db.WriteChatInfoToLocalFile(chatRoom)
	chat_db.SaveRoomNameListToFile(chatRoomList)
	utils.WriteResponse(w, 0, "Success")
}
func HistoryMessagesHandler(w http.ResponseWriter, r *http.Request) {
	// 确保关闭请求体
	defer r.Body.Close()
	// 检查请求方法是否为POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
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
	if !ChatRoomExist(roomName) {
		utils.WriteResponse(w, chat_type.ErrorInvalidInput, "Chat room not exist")
		return
	}
	chatRoom := getChatRoomByName(roomName)
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
		slog.Error("Failed to convert message to JSON", "error", err)
		return
	}
	slog.Info("HistoryMessagesHandler", "messages size", len(messages))
	utils.WriteResponseWithData(w, chat_type.ErrorCodeSuccess, "Success", jsonMsg)
}
func ChatRoomExist(chatRoomName string) bool {
	for _, room := range chatRoomList {
		if room == chatRoomName {
			return true
		}
	}
	return false
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
	_, err = w.Write(jsonMsg)
	if err != nil {
		slog.Error("Failed to write response", "error", err)
		return
	}
}

func InitChatRoom() {
	chatRoomList = chat_db.LoadRoomNameListFromFile()
	for _, room := range chatRoomList {
		chatRoom := getChatRoomByName(room)
		chatRoomMap[room] = chatRoom
	}
}

func ChatRoomHandler(w http.ResponseWriter, r *http.Request) {
	chatRoomList = chat_db.LoadRoomNameListFromFile()
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade connection to WebSocket", http.StatusInternalServerError)
		return
	}

	// 获取request中的参数,比如id和chatroom
	id := r.URL.Query().Get("id")
	chatRoomName := r.URL.Query().Get("chatroom")

	// Get chat room by name
	if chatRoomName == "" || chatRoomName == "null" {
		chatRoomName = "default"
	}

	go handleClientConn(conn, id, chatRoomName)
}

func handleClientConn(conn *websocket.Conn, id, chatRoomName string) {
	defer conn.Close()
	chatRoom := getChatRoomByName(chatRoomName)
	chat_db.SaveRoomNameListToFile(chatRoomList)
	// Add user to chat room
	addNewUserToChatRoom(chatRoom, id)
	// Add connection to the list of active connections
	chatRoom.Connections = append(chatRoom.Connections, conn)
	slog.Info("new user join", "id", id, "roomName", chatRoomName, "memberSize", len(chatRoom.Connections))
	sendMessage(chat_type.Message{Type: "roomList", ChatRoomList: chatRoomList, RoomName: chatRoomName}, conn)
	// Read messages from the WebSocket connection
	for {
		// Read message from the WebSocket
		_, msg, err := conn.ReadMessage()
		if err != nil { // remove connection from the list of active connections
			slog.Error("Failed to read message from WebSocket", "id", id, "error", err)
			chatRoom.Connections = removeConnection(conn, chatRoom.Connections)
			slog.Info("Leave, chatRoom.Connections", "id", id, "roomName", chatRoomName, "memberSize", len(chatRoom.Connections))
			if len(chatRoom.Connections) == 0 {
				CloseChatRoom(chatRoom)
			}
			break
		}
		// Parse message into Message struct
		var message chat_type.Message
		err = json.Unmarshal(msg, &message)
		if err != nil {
			slog.Error("Failed to parse message", "error", err)
			continue
		}
		// 限制message.Image的文件大小
		if message.Image != "" {
			if len(message.Image) > 1024*1024*20 {
				message.Image = ""
			}
		}
		utils.TryTransferImagePathToMessage(&message)
		// get tUer
		//tUer := user.GetUserById(id)
		// 暂时直接使用消息带上来的userName
		//message.UserName = tUer.UserName
		message.MsgID = utils.GenerateId()
		message.SendTime = utils.GetCurTime()
		// fmt.Println("message:", message)
		// Add message to the list of all messages
		chatRoom.Messages = append(chatRoom.Messages, message)
		chat_db.WriteChatInfoToLocalFile(chatRoom)
		// Broadcast message to all active connections
		for _, c := range chatRoom.Connections {
			err := sendMessage(message, c)
			if err != nil {
				// Remove connection from the list of active connections
				for i, conn := range chatRoom.Connections {
					if conn == c {
						chatRoom.Connections = append(chatRoom.Connections[:i], chatRoom.Connections[i+1:]...)
						break
					}
				}
				break
			}
		}
	}
}
