package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
	chat_type "web_server/type"
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

func TryTransferImagePathToMessage(message *chat_type.Message) error {
	if message.Image != "" {
		message.Type = "image"
		fmt.Println("start to decode image data")
		fragments := strings.Split(message.Image, ",")
		if len(fragments) > 1 {
			b64data := fragments[1]
			imageData, err := base64.StdEncoding.DecodeString(b64data)
			if err != nil {
				fmt.Println("Failed to decode image data:", err)
				return err
			}
			hasher := md5.New()
			hasher.Write(imageData)
			imageFileName := hex.EncodeToString(hasher.Sum(nil)) + ".jpg"
			err = ioutil.WriteFile(ImageDir+imageFileName, imageData, 0644)
			if err != nil {
				fmt.Println("Failed to write image file:", err)
				return err
			}
			message.Image = ImageDir + imageFileName
			fmt.Println("Image saved to:", message.Image)
		} else {
			u, err := url.Parse(fragments[0])
			if err != nil {
				fmt.Println("Failed to parse image url:", err)
				return err
			}
			message.Image = ImageDir + u.String()
			fmt.Println("Image url from:", message.Image)
		}
	} else {
		message.Type = "text"
	}
	return nil
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
