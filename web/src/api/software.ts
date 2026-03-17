import request from '@/utils/request'

export interface SoftwareTask {
  id: string
  store_id: string
  operator: string
  file_name: string
  file_size: number
  install_args: string
  status: 'pending' | 'uploading' | 'downloading' | 'installing' | 'installed' | 'failed'
  message: string
  created_at: string
  completed_at: string | null
}

export interface PushSoftwareRequest {
  file_name: string
  file_url: string
  file_size: number
  install_args?: string
}

export function pushSoftware(storeId: string, data: PushSoftwareRequest): Promise<{ task_id: string }> {
  return request.post(`/stores/${storeId}/software`, data)
}

export function getSoftwareStatus(storeId: string, taskId: string): Promise<SoftwareTask> {
  return request.get(`/stores/${storeId}/software/${taskId}`)
}