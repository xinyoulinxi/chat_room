package utils

import (
	"encoding/json"
	"log/slog"
	"os"
	"strconv"
)

const (
	DataPath    = "data/"
	FileDir     = "data/files/"
	ImageDir    = "data/images/"
	defaultPort = "80"
)

type (
	config struct {
		IP   string `json:"ipAddress"`
		Port int    `json:"port"`
	}
)

func GetLocalIP() string {
	// 读取配置文件, 获取IP地址，配置文件格式如下：
	// {
	//     "ipAddress": 0.0.0.0",
	//      "Port":80
	// }
	filePath := "./config.json"
	data, err := os.ReadFile(filePath)
	if err != nil {
		// 打印英文日志：读取配置文件失败
		slog.Warn("Read config file failed, use default config", "error", err)
		return ":" + defaultPort
	}
	config := new(config)
	err = json.Unmarshal(data, &config)
	if err != nil {
		slog.Warn("Parse config file failed, use default config", "error", err)
		return ":" + defaultPort
	}
	if config.Port <= 0 || config.Port > 65535 {
		// 翻译：端口范围只能在1~65535
		slog.Warn("Port range must be between 1 and 65535, use default port", "error", err)
		return config.IP + ":" + defaultPort
	}
	return config.IP + ":" + strconv.Itoa(config.Port)
}

func InitEnv() {
	dirs := []string{FileDir, ImageDir}
	for _, dir := range dirs {
		if err := EnsureDir(dir); err != nil {
			panic(err)
		}
	}

}
