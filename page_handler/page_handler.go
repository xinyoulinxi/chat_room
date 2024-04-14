package page_handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	chat_room "web_server/chat_room"
	utils "web_server/utils"

	"github.com/gorilla/handlers"
	// "github.com/gorilla/handlers"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "pages/index.html")
}

func sumHandler(w http.ResponseWriter, r *http.Request) {
	// 解析查询参数
	num1Str := r.URL.Query().Get("num1")
	num2Str := r.URL.Query().Get("num2")

	// 将字符串参数转换为整数
	num1, err1 := strconv.Atoi(num1Str)
	num2, err2 := strconv.Atoi(num2Str)

	// 错误处理
	if err1 != nil || err2 != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// 计算和
	sum := num1 + num2

	// 将结果转换为JSON并返回
	result := map[string]int{"sum": sum}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

type noStore struct {
	h http.Handler
}

func (n *noStore) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	n.h.ServeHTTP(w, r)
}

func StartWebServer() {
	url := utils.GetLocalIP()
	mux := http.NewServeMux()

	// 配置CORS
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost, http.MethodOptions})
	// 创建原始的文件服务器
	fs := http.FileServer(http.Dir("pages"))

	// 创建包装器
	noStoreFs := &noStore{h: fs}

	// 使用包装器替代原始的文件服务器
	mux.Handle("/pages/", http.StripPrefix("/pages/", noStoreFs))
	// 主页
	mux.Handle("/", utils.NoCacheMiddleware(http.HandlerFunc(indexHandler)))
	// Image file service
	imageFs := http.FileServer(http.Dir("data"))
	mux.Handle("/data/", http.StripPrefix("/data/", imageFs))

	// Page 1
	mux.Handle("/page_1", utils.NoCacheMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pages/page_1.html")
	})))

	// login
	mux.Handle("/login", utils.NoCacheMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("login")
		http.ServeFile(w, r, "pages/login.html")
	})))

	// Chat Room
	mux.Handle("/chat_room", utils.NoCacheMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("chat_room")
		http.ServeFile(w, r, "pages/chat_room.html")
	})))

	// sum
	mux.HandleFunc("/sum", sumHandler)

	// room list
	mux.HandleFunc("/room_list", chat_room.ChatRoomListHandler)
	// create Chat Room
	mux.HandleFunc("/create_room", chat_room.CreateChatRoomHandler)

	// ws
	mux.HandleFunc("/ws", chat_room.ChatRoomHandler)

	fmt.Printf("Starting server at http://%s\n", url+":8080/")
	fmt.Printf("Starting server at http://%s\n", url+":8080/page_1")
	fmt.Printf("Starting server at http://%s\n", url+":8080/login")
	err := http.ListenAndServe(url+":8080", handlers.CORS(headersOk, originsOk, methodsOk)(mux))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
