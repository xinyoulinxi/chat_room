package utils

import (
	"encoding/base64"
	"log/slog"
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
	return time.Now().Format(time.DateTime)
}
func GetRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// saveFile 根据Base64编码的数据和推断的文件类型保存文件
func saveFile(message *chat_type.Message) error {
	fragments := strings.Split(message.Image, ",")
	base64Data := fragments[1]
	name := message.Content
	// 解码Base64字符串
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		slog.Error("Failed to decode base64 data", "error", err)
		return err
	}

	// 推断文件类型
	kind, err := filetype.Match(data)
	if err != nil {
		slog.Error("Failed to infer file type", "error", err)
		return err
	}

	// 获取对应的文件扩展名
	ext := kind.Extension

	// 判断文件类型，如果是图片类型，则给message.Type赋值为image
	slog.Info("Detected file type", "type", kind.MIME.Type)

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
	err = os.WriteFile(filename, data, 0666)
	if err != nil {
		return err
	}

	return nil
}

func TryTransferImagePathToMessage(message *chat_type.Message) {
	slog.Info("TryTransferImagePathToMessage")
	if message.Image != "" {
		message.Type = "image"
		slog.Info("start to decode image data")
		fragments := strings.Split(message.Image, ",")
		if len(fragments) > 1 {
			if err := saveFile(message); err != nil {
				slog.Error("Failed to save file", "error", err)
			}
		} else {
			u, err := url.Parse(fragments[0])
			if err != nil {
				slog.Error("Failed to parse image url", "error", err)
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
func NormalCacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "Cache-Control: public, max-age=31536000") // HTTP 1.1.
		next.ServeHTTP(w, r)
	})
}
func EnsureDir(dirName string) error {
	_, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dirName, 0755)
		if err != nil {
			slog.Error("check directory failed", "dir", dirName, "error", err)
			return err
		}
	}
	return nil
}

func EnsureFileExist(filePath string) error {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		err = os.WriteFile(filePath, []byte{}, 0644)
		if err != nil {
			slog.Error("check file failed", "file", filePath, "error", err)
			return err
		}
	}
	return nil
}
