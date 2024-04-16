package utils

import (
	"encoding/json"
	"net/http"
	chat_type "web_server/type"
)

func GetReturnMessageJson(errorCode int, message string) []byte {
	jsonMsg, _ := json.Marshal(chat_type.ReturnMessage{ErrorCode: errorCode, Message: message})
	return jsonMsg
}

func WriteResponse(w http.ResponseWriter, code int, message string) {
	_, err := w.Write(GetReturnMessageJson(code, message))
	if err != nil {
		return
	}
}
