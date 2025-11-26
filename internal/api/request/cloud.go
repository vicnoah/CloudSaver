package request

// GetShareInfoRequest 获取分享信息请求
type GetShareInfoRequest struct {
	ShareCode string `form:"shareCode" binding:"required"`
	Passcode  string `form:"passcode"`
}

// GetFolderListRequest 获取文件夹列表请求
type GetFolderListRequest struct {
	ParentCid string `form:"parentCid"`
}

// SaveFileRequest 保存文件请求
type SaveFileRequest struct {
	ShareCode  string   `json:"shareCode" binding:"required"`
	Passcode   string   `json:"passcode"`
	FolderID   string   `json:"folderId" binding:"required"`
	FileIDs    []string `json:"fileIds" binding:"required"`
	FileTokens []string `json:"fileTokens,omitempty"` // 夸克专用
	PwdID      string   `json:"pwdId,omitempty"`      // 夸克专用
	Stoken     string   `json:"stoken,omitempty"`     // 夸克专用
}
