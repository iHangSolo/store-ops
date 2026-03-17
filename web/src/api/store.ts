import request from '@/utils/request'
import type { PageResponse } from '@/utils/request'

export interface Store {
  id: string
  name: string
  device_id: string
  device_fingerprint: string
  token: string
  rustdesk_id: string
  client_version: string
  status: 'pending' | 'approved' | 'rejected'
  is_online: boolean
  cpu_percent: number
  memory_percent: number
  last_seen: string | null
  created_at: string
  updated_at: string
}

export interface StoreListParams {
  status?: string
  page?: number
  page_size?: number
}

export function getStores(params: StoreListParams): Promise<PageResponse<Store>> {
  return request.get('/stores', { params })
}

export function getStore(id: string): Promise<Store> {
  return request.get(`/stores/${id}`)
}

export function approveStore(id: string): Promise<void> {
  return request.post(`/stores/${id}/approve`)
}

export function rejectStore(id: string, reason: string): Promise<void> {
  return request.post(`/stores/${id}/reject`, { reason })
}

export function deleteStore(id: string): Promise<void> {
  return request.delete(`/stores/${id}`)
}

export function regenerateToken(id: string): Promise<{ token: string }> {
  return request.post(`/stores/${id}/regenerate-token`)
}