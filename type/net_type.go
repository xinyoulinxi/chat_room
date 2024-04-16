package chat_type

type ReturnMessage struct {
	ErrorCode int    `json:"errorCode"`
	Message   string `json:"message"`
}

const (
	ErrorCodeSuccess     = 0
	ErrorCodeFail        = 1
	ErrorInvalidInput    = 2
	ErrorUserExist       = 3
	ErrorUserNotExist    = 4
	ErrorPasswordError   = 5
	ErrorInvalidPassword = 6
)
