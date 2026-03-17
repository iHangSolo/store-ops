import request from '@/utils/request'

export interface DesktopStatus {
  store_id: string
  rustdesk_id: string
  is_online: boolean
  has_active_session: boolean
  queue_count: number
  current_session?: {
    id: string
    operator: string
    started_at: string
  }
}

export interface QueuePosition {
  position: number
  estimated_wait_seconds: number
}

export function getDesktopStatus(storeId: string): Promise<DesktopStatus> {
  return request.get(`/stores/${storeId}/desktop/status`)
}

export function connectDesktop(storeId: string): Promise<{ rustdesk_url: string }> {
  return request.post(`/stores/${storeId}/desktop/connect`)
}

export function queueDesktop(storeId: string): Promise<QueuePosition> {
  return request.post(`/stores/${storeId}/desktop/queue`)
}

export function endSession(storeId: string, sessionId: string): Promise<void> {
  return request.delete(`/stores/${storeId}/desktop/session/${sessionId}`)
}

export function leaveQueue(storeId: string, queueId: string): Promise<void> {
  return request.delete(`/stores/${storeId}/desktop/queue/${queueId}`)
}