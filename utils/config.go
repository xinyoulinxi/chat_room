package utils

import (
	"encoding/json"
	"io/ioutil"
)

const DataPath = "data/"
const ChatRoomPath = "data/chatroom/"
const FileDir = "data/files/"
const ImageDir = "data/images/"
const RoomListPath = "data/room_list.json"

func GetLocalIP() string {
	// 读取配置文件, 获取IP地址，配置文件格式如下：
	// {
	//     "ipAddress": 0.0.0.0"
	// }
	filePath := "./config.json"
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return ""
	}

	var config struct {
		IP string `json:"ipAddress"`
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return ""
	}

	return config.IP
}
func EnsureDirEnv() {
	EnsureDir(DataPath)
	EnsureDir(FileDir)
	EnsureDir(ImageDir)
	EnsureDir(ChatRoomPath)
	EnsureFileExist(RoomListPath)

}
func InitEnv() {
	EnsureDirEnv()
}
