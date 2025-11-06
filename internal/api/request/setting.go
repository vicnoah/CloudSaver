package request

// SaveSettingRequest 保存设置请求
type SaveSettingRequest struct {
	Cloud115Cookie string `json:"cloud115Cookie"`
	QuarkCookie    string `json:"quarkCookie"`
}
