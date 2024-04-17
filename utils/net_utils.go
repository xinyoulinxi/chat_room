package utils

import (
	"encoding/json"
	"log/slog"
	"net/http"
	chat_type "web_server/type"
)

func GetReturnMessageJson(errorCode int, message string) []byte {
	jsonMsg, _ := json.Marshal(chat_type.ReturnMessage{ErrorCode: errorCode, Message: message})
	return jsonMsg
}

func GetReturnMessageJsonWithData(errorCode int, message string, data []byte) []byte {
	jsonMsg, _ := json.Marshal(chat_type.ReturnMessage{ErrorCode: errorCode, Message: message, Data: data})
	return jsonMsg
}

func WriteResponseWithData(w http.ResponseWriter, code int, message string, data []byte) {
	slog.Info("WriteResponse", "code", code, "message", message)
	_, err := w.Write(GetReturnMessageJsonWithData(code, message, data))
	if err != nil {
		return
	}
}

func WriteResponse(w http.ResponseWriter, code int, message string) {
	slog.Info("WriteResponse", "code", code, "message", message)
	_, err := w.Write(GetReturnMessageJson(code, message))
	if err != nil {
		return
	}
}
