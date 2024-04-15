package chat_type

import "encoding/json"

type ReturnMessage struct {
	ErrorCode int    `json:"errorCode"`
	Message   string `json:"message"`
}

func GetReturnMessageJson(errorCode int, message string) []byte {
	jsonMsg, _ := json.Marshal(ReturnMessage{ErrorCode: errorCode, Message: message})
	return jsonMsg
}
