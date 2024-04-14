package chat_room

import (
	"encoding/json"
	"fmt"
	"net/http"
	chat_db "web_server/db"
	chat_type "web_server/type"
	utils "web_server/utils"

	"github.com/gorilla/websocket"
	// "github.com/gorilla/handlers"
)

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
			fmt.Println("add new chat room:", roomName)
			fmt.Println(chatRoomList)
			chatRoomList = append(chatRoomList, roomName)
			fmt.Println(chatRoomList)
		}
		return &chatRoom
	}
}

func sendChatRoomMessagesToNewUser(chatRoom *chat_type.ChatRoom, conn *websocket.Conn) {
	// Send all messages to the newly connected user
	jsonMsg, err := json.Marshal(chatRoom.Messages)
	if err != nil {
		fmt.Println("Failed to convert message to JSON:", err)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, jsonMsg)
	if err != nil {
		fmt.Println("Failed to send message to user:", err)
		return
	}
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
	jsonMsg, err := json.Marshal([]chat_type.Message{message})
	fmt.Println("send msg:", string(jsonMsg))
	if err != nil {
		fmt.Println("Failed to convert message to JSON:", err)
		return err
	}

	// Write message to WebSocket
	err = c.WriteMessage(websocket.TextMessage, jsonMsg)
	return err
}

func CreateChatRoomHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request parameters
	roomName := r.URL.Query().Get("roomName")
	fmt.Println("CreateChatRoomHandler", roomName)
	if roomName == "" {
		w.Write(chat_type.GetReturnMessageJson(1, "Invalid chat room name"))
		return
	}

	for _, room := range chatRoomList {
		if room == roomName {
			w.Write(chat_type.GetReturnMessageJson(1, "Chat room already exists"))
			return
		}
	}

	chatRoom := getChatRoomByName(roomName)
	chat_db.WriteChatInfoToLocalFile(chatRoom)
	chat_db.SaveRoomNameListToFile(chatRoomList)
	w.Write(chat_type.GetReturnMessageJson(0, "Create chat room successfully"))
}

func ChatRoomHandler(w http.ResponseWriter, r *http.Request) {
	chatRoomList = chat_db.LoadRoomNameListFromFile()
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade connection to WebSocket", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	// 获取request中的参数,比如id和chatroom
	id := r.URL.Query().Get("id")
	chatRoomName := r.URL.Query().Get("chatroom")
	fmt.Println("new user join id:", id, ", Room Name:", chatRoomName)
	// Get chat room by name
	if chatRoomName == "" || chatRoomName == "null" {
		chatRoomName = "default"
	}
	chatRoom := getChatRoomByName(chatRoomName)
	chat_db.SaveRoomNameListToFile(chatRoomList)
	// Add user to chat room
	addNewUserToChatRoom(chatRoom, id)
	// Add connection to the list of active connections
	chatRoom.Connections = append(chatRoom.Connections, conn)
	fmt.Println("new connection current room user=", len(chatRoom.Connections))
	sendMessage(chat_type.Message{Type: "roomList", ChatRoomList: chatRoomList, RoomName: chatRoomName}, conn)
	sendChatRoomMessagesToNewUser(chatRoom, conn)
	// Read messages from the WebSocket connection
	for {
		// Read message from the WebSocket
		_, msg, err := conn.ReadMessage()
		if err != nil { // remove connection from the list of active connections
			fmt.Println("Failed to connect :", err)
			chatRoom.Connections = removeConnection(conn, chatRoom.Connections)
			fmt.Println(id, "Leave, chatRoom.Connections:", len(chatRoom.Connections))
			if len(chatRoom.Connections) == 0 {
				CloseChatRoom(chatRoom)
			}
			break
		}
		// Parse message into Message struct
		var message chat_type.Message
		err = json.Unmarshal(msg, &message)
		if err != nil {
			fmt.Println("Failed to parse message:", err)
			continue
		}
		utils.TryTransferImagePathToMessage(&message)
		message.SendTime = utils.GetCurTime()
		fmt.Println("message:", message)
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
