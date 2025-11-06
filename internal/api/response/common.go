package response

// Response 通用响应结构
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
	Code    string      `json:"code,omitempty"`
}

// Success 成功响应
func Success(data interface{}) Response {
	return Response{
		Success: true,
		Data:    data,
	}
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(data interface{}, message string) Response {
	return Response{
		Success: true,
		Data:    data,
		Message: message,
	}
}

// Error 错误响应
func Error(message string) Response {
	return Response{
		Success: false,
		Error:   message,
	}
}

// ErrorWithCode 带错误码的错误响应
func ErrorWithCode(code, message string) Response {
	return Response{
		Success: false,
		Error:   message,
		Code:    code,
	}
}

// Pagination 分页信息
type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

// PaginatedData 分页数据
type PaginatedData struct {
	Items      interface{} `json:"items"`
	Pagination Pagination  `json:"pagination"`
}

// Paginated 分页响应
func Paginated(items interface{}, page, pageSize, total int) Response {
	return Response{
		Success: true,
		Data: PaginatedData{
			Items: items,
			Pagination: Pagination{
				Page:     page,
				PageSize: pageSize,
				Total:    total,
			},
		},
	}
}
