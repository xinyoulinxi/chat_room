package utils

import (
	"encoding/base64"
	"github.com/h2non/filetype"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GetRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// SaveFile 根据Base64编码的数据和推断的文件类型保存文件
func SaveFile(name string, fileData *string) (string, string, error) {
	fragments := strings.Split(*fileData, ",")
	base64Data := fragments[1]
	// 解码Base64字符串
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		slog.Error("Failed to decode base64 data", "error", err)
		return "", "", err
	}

	// 推断文件类型
	kind, err := filetype.Match(data)
	if err != nil {
		slog.Error("Failed to infer file type", "error", err)
		return "", "", err
	}

	// 获取对应的文件扩展名
	ext := kind.Extension

	// 判断文件类型，如果是图片类型，则给message.Type赋值为image
	slog.Info("Detected file type", "type", kind.MIME.Type)

	filename := ""
	fileType := ""
	if kind.MIME.Type == "image" {
		filename = ImageDir + name
		fileType = "image"
	} else {
		filename = FileDir + name
		fileType = "file"
	}
	if ext == "" {
		ext = "bin" // 使用通用二进制扩展名作为后备
	}

	// 保存文件
	err = os.WriteFile(filename, data, 0666)
	if err != nil {
		return "", fileType, err
	}

	return filename, fileType, nil
}

// EnsureDir 确保目录存在
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

// EnsureFileExist 确保文件存在
func EnsureFileExist(filePath string) error {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		if err = EnsureDir(filepath.Dir(filePath)); err != nil {
			return err
		}
		err = os.WriteFile(filePath, []byte{}, 0644)
		if err != nil {
			slog.Error("check file failed", "file", filePath, "error", err)
			return err
		}
	}
	return nil
}
