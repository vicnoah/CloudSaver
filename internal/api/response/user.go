package response

// UserLoginResponse 用户登录响应
type UserLoginResponse struct {
	User  UserInfo `json:"user"`
	Token string   `json:"token"`
}

// UserInfo 用户信息
type UserInfo struct {
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	Role     int    `json:"role"`
}
