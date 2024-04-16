package user

import (
	"encoding/json"
	"log/slog"
	"net/http"
	chat_db "web_server/db"
	chat_type "web_server/type"
	"web_server/utils"
)

var userList = make([]chat_type.User, 0)
var userMap = make(map[string]chat_type.User)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// 确保关闭请求体
	defer r.Body.Close()
	// 检查请求方法是否为POST
	if r.Method != http.MethodPost {
		utils.WriteResponse(w, chat_type.ErrorCodeFail, "Invalid request method")
		return
	}

	// 解析请求体
	var registerData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&registerData)
	if err != nil {
		utils.WriteResponse(w, chat_type.ErrorCodeFail, "Failed to parse request body")
		return
	}
	userName := registerData.Username
	passWord := registerData.Password
	slog.Info("RegisterHandler", "username", userName, "password", passWord)
	if userName == "" || passWord == "" {
		utils.WriteResponse(w, chat_type.ErrorInvalidInput, "Invalid input")
		return
	}
	if UserExist(userName) {
		utils.WriteResponse(w, chat_type.ErrorUserExist, "User already exist")
		return
	}
	user := addUser(userName, passWord)
	slog.Info("RegisterHandler", "user", user)
	if user != nil {
		utils.WriteResponse(w, chat_type.ErrorCodeSuccess, user.UserID)
	} else {
		utils.WriteResponse(w, chat_type.ErrorCodeFail, "Failed to add user")
	}
}

func InitUserInfos() {
	updateUserInfosFromLocalFile()
}

func updateUserInfosFromLocalFile() {
	userList = chat_db.LoadUsersFromLocalFile()
	for _, user := range userList {
		userMap[user.UserID] = user
	}
}

func saveUserInfosToLocalFile() {
	err := chat_db.WriteUsersToLocalFile(userList)
	if err != nil {
		return
	}
}
func getUserByName(userName string) *chat_type.User {
	for _, user := range userList {
		if user.UserName == userName {
			return &user
		}
	}
	return nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// 确保关闭请求体
	defer r.Body.Close()
	// 检查请求方法是否为POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求体
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		utils.WriteResponse(w, chat_type.ErrorCodeFail, "Failed to parse request body")
		return
	}
	userName := loginData.Username
	passWord := loginData.Password

	slog.Info("LoginHandler", "username", userName, "password", passWord)
	if userName == "" || passWord == "" {
		utils.WriteResponse(w, chat_type.ErrorInvalidInput, "Invalid input")
		return
	}

	if !UserExist(userName) {
		utils.WriteResponse(w, chat_type.ErrorUserNotExist, "User not exist")
		return
	}

	if !CheckPassword(userName, passWord) {
		utils.WriteResponse(w, chat_type.ErrorInvalidPassword, "Invalid password")
		return
	}
	// 返回整个user结构体
	user := getUserByName(userName)
	if user == nil {
		utils.WriteResponse(w, chat_type.ErrorCodeFail, "Failed to get user")
		return
	}
	utils.WriteResponse(w, chat_type.ErrorCodeSuccess, user.UserID)
}

func UserRegisted(userId string) bool {
	for _, user := range userList {
		if user.UserID == userId {
			return true
		}
	}
	return false
}

func UserExist(userName string) bool {
	for _, user := range userList {
		if user.UserName == userName {
			return true
		}
	}
	return false
}

func CheckPassword(userName string, passWord string) bool {
	for _, user := range userList {
		if user.UserName == userName && user.PassWord == passWord {
			return true
		}
	}
	return false
}

// 产生一个不会重复的用户ID
func createUserId() string {
	// 生成一个随机的用户ID
	userId := utils.GetRandomString(16)
	// 确保用户ID不会重复
	for {
		if _, ok := userMap[userId]; ok {
			userId = utils.GetRandomString(16)
		} else {
			break
		}
	}
	return userId
}

func addUser(userName string, passWord string) *chat_type.User {
	if userName == "" || passWord == "" {
		return nil
	}
	userId := createUserId()
	slog.Info("addUser", "userName", userName, "passWord", passWord, "userId", userId)
	user := chat_type.User{UserID: userId, UserName: userName, PassWord: passWord}
	userList = append(userList, user)
	userMap[user.UserID] = user
	saveUserInfosToLocalFile()
	return &user
}

func GetUserById(id string) *chat_type.User {
	if user, ok := userMap[id]; ok {
		return &user
	}
	return nil
}
