import apiClient from '../client'

export interface DoubanParams {
  type: string
  tag: string
  page_limit?: string
  page_start?: string
}

export interface DoubanItem {
  id: string
  title: string
  rate: string
  cover: string
  url: string
  is_new: boolean
}

export const doubanApi = {
  getHotList(params: DoubanParams) {
    return apiClient.get<DoubanItem[]>('/douban/hot', { params })
  },
}
