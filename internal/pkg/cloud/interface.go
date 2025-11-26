package cloud

import (
	"context"
	"cloudsaver/internal/api/response"
)

// CloudStorageService 云盘存储服务接口
type CloudStorageService interface {
	// SetCookie 设置Cookie
	SetCookie(userUUID string) error

	// GetShareInfo 获取分享信息
	GetShareInfo(ctx context.Context, shareCode string, passcode string) (*response.ShareInfoResponse, error)

	// GetFolderList 获取文件夹列表
	GetFolderList(ctx context.Context, parentID string) (*response.FolderListResponse, error)

	// SaveSharedFile 保存分享文件
	SaveSharedFile(ctx context.Context, shareCode, passcode, folderID string, fileIDs []string, fileTokens []string, pwdID, stoken string) error
}
