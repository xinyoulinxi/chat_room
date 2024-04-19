package chat_type

type User struct {
	UserID   string `json:"userId"`
	UserName string `json:"userName"`
	PassWord string `json:"passWord"`
	Avatar   string `json:"avatar"`
}

type Users []*User
