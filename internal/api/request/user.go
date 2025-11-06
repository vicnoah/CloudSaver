package request

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username     string `json:"username" binding:"required,min=3,max=20"`
	Password     string `json:"password" binding:"required,min=6,max=32"`
	RegisterCode int    `json:"registerCode" binding:"required"`
}
