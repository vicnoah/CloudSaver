import apiClient from '../client'
import type { ShareInfoResponse, FolderListResponse, SaveFileParams } from '@/types/api'

export const cloud115Api = {
  getShareInfo(shareCode: string, passcode?: string) {
    return apiClient.get<ShareInfoResponse>('/cloud115/share-info', {
      params: { shareCode, passcode }
    })
  },
  
  getFolderList(parentCid?: string) {
    return apiClient.get<FolderListResponse>('/cloud115/folders', {
      params: { parentCid }
    })
  },
  
  saveFile(params: SaveFileParams) {
    return apiClient.post('/cloud115/save', params)
  },
}
