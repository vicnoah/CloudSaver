import apiClient from '../client'

export interface SearchParams {
  keyword?: string
  channelId?: string
  messageId?: string
}

export interface SearchResult {
  id: string
  channelInfo: {
    id: string
    name: string
    channelLogo: string
  }
  list: Array<{
    messageId: string
    title: string
    content: string
    pubDate: string
    image: string
    cloudLinks: string[]
    cloudType: string
    tags: string[]
    channel: string
    channelId: string
  }>
}

export const searchApi = {
  search(params: SearchParams) {
    return apiClient.get<SearchResult[]>('/search', { params })
  },
}
