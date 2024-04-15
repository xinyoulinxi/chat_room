package utils

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	chat_type "web_server/type"

	"github.com/h2non/filetype"
)

func GetChatRoomFilePath(chatName string) string {
	if chatName == "" {
		return ChatRoomPath + "chatroom.json"
	}
	return ChatRoomPath + "chatroom_" + chatName + ".json"
}

func GetCurTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// saveFile 根据Base64编码的数据和推断的文件类型保存文件
func saveFile(message *chat_type.Message) error {
	fmt.Println("saveFile")

	fragments := strings.Split(message.Image, ",")
	base64Data := fragments[1]
	name := message.Content
	// 解码Base64字符串
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		fmt.Println("Failed to decode base64 data:", err)
		return err
	}

	// 推断文件类型
	kind, err := filetype.Match(data)
	if err != nil {
		fmt.Println("Failed to infer file type:", err)
		return err
	}

	// 获取对应的文件扩展名
	ext := kind.Extension

	// 判断文件类型，如果是图片类型，则给message.Type赋值为image
	fmt.Println("kind.MIME.Type", kind.MIME.Type)
	filename := ""
	if kind.MIME.Type == "image" {
		filename = ImageDir + name
		message.Type = "image"
		message.Image = filename
	} else {
		filename = FileDir + name
		message.Type = "file"
		message.File = filename
	}
	if ext == "" {
		ext = "bin" // 使用通用二进制扩展名作为后备
	}

	// 保存文件
	err = ioutil.WriteFile(filename, data, 0666)
	if err != nil {
		return err
	}

	return nil
}

func TryTransferImagePathToMessage(message *chat_type.Message) {
	fmt.Println("TryTransferImagePathToMessage")
	if message.Image != "" {
		message.Type = "image"
		fmt.Println("start to decode image data")
		fragments := strings.Split(message.Image, ",")
		if len(fragments) > 1 {
			saveFile(message)
		} else {
			u, err := url.Parse(fragments[0])
			if err != nil {
				fmt.Println("Failed to parse image url:", err)
				return
			}
			message.Image = u.String()
		}
	} else {
		message.Type = "text"
	}
}

// noCacheMiddleware 为响应添加缓存控制头
func NoCacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
		w.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
		w.Header().Set("Expires", "0")                                         // Proxies.
		next.ServeHTTP(w, r)
	})
}

func EnsureDir(dirName string) error {
	_, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dirName, 0755)
		if errDir != nil {
			return err
		}
	}
	return nil
}

func EnsureFileExist(filePath string) error {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		errFile :=
			ioutil.WriteFile(filePath, []byte{}, 0644)
		if errFile != nil {
			return errFile
		}
	}
	return nil
}
