import apiClient from '../client'
import type { ShareInfoResponse, FolderListResponse, SaveQuarkFileParams } from '@/types/api'

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
  
  saveFile(params: SaveQuarkFileParams) {
    return apiClient.post('/quark/save', params)
  },
}
