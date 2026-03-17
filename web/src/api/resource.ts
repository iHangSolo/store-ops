import request from '@/utils/request'
import type { PageResponse } from '@/utils/request'

export interface ResourceInfo {
  cpu_percent: number
  memory_percent: number
  memory_available_gb: number
  disks: DiskInfo[]
  recorded_at: string
}

export interface DiskInfo {
  name: string
  percent: number
  total_gb: number
  available_gb: number
}

export interface ResourceHistoryParams {
  start_time?: string
  end_time?: string
  page?: number
  page_size?: number
}

export function getResource(storeId: string): Promise<ResourceInfo> {
  return request.get(`/stores/${storeId}/resource`)
}

export function getResourceHistory(storeId: string, params: ResourceHistoryParams): Promise<PageResponse<ResourceInfo>> {
  return request.get(`/stores/${storeId}/resource/history`, { params })
}