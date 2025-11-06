// API 响应类型
export interface ApiResponse<T = any> {
  success: boolean
  data?: T
  error?: string
  message?: string
  code?: string
}

// 用户相关类型
export interface UserInfo {
  uuid: string
  username: string
  role: number
}

export interface LoginParams {
  username: string
  password: string
}

export interface RegisterParams {
  username: string
  password: string
}

export interface LoginResponse {
  user: UserInfo
  token: string
}

// 设置相关类型
export interface UserSetting {
  cloud115Cookie: string
  quarkCookie: string
}

// 云盘相关类型
export interface FileItem {
  file_id: string
  file_name: string
  file_type: number
  file_id_token?: string
}

export interface ShareInfoResponse {
  list: FileItem[]
  pwd_id?: string
  stoken?: string
  file_size: number
}

export interface FolderItem {
  cid: string
  name: string
  path?: FolderItem[]
}

export interface FolderListResponse {
  folders: FolderItem[]
}

export interface SaveFileParams {
  shareCode: string
  passcode?: string
  folderId: string
  fileIds: string[]
  fileTokens?: string[]
  pwdId?: string
  stoken?: string
}
