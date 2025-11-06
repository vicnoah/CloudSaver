import apiClient from '../client'
import type { ShareInfoResponse, FolderListResponse, SaveFileParams } from '@/types/api'

export const quarkApi = {
  getShareInfo(shareCode: string, passcode?: string) {
    return apiClient.get<ShareInfoResponse>('/quark/share-info', {
      params: { shareCode, passcode }
    })
  },
  
  getFolderList(parentCid?: string) {
    return apiClient.get<FolderListResponse>('/quark/folders', {
      params: { parentCid }
    })
  },
  
  saveFile(params: SaveFileParams) {
    return apiClient.post('/quark/save', params)
  },
}
