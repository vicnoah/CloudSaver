package response

// ShareInfoResponse 分享信息响应
type ShareInfoResponse struct {
	List     []FileItem `json:"list"`
	PwdID    string     `json:"pwd_id,omitempty"`     // 夸克专用
	Stoken   string     `json:"stoken,omitempty"`     // 夸克专用
	FileSize int64      `json:"file_size"`
}

// FileItem 文件项
type FileItem struct {
	FileID      string `json:"file_id"`
	FileName    string `json:"file_name"`
	FileType    int    `json:"file_type"`                // 1-文件夹 0-文件
	FileIDToken string `json:"file_id_token,omitempty"` // 夸克分享令牌
}

// FolderListResponse 文件夹列表响应
type FolderListResponse struct {
	Folders []FolderItem `json:"folders"`
}

// FolderItem 文件夹项
type FolderItem struct {
	CID  string       `json:"cid"`
	Name string       `json:"name"`
	Path []FolderItem `json:"path,omitempty"`
}
