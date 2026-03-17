import request from '@/utils/request'
import type { PageResponse } from '@/utils/request'

export interface AuditLog {
  id: string
  store_id: string | null
  store_name: string
  action_type: string
  action: string
  detail: string
  operator: string
  created_at: string
}

export interface AuditLogParams {
  store_id?: string
  action_type?: string
  start_date?: string
  end_date?: string
  page?: number
  page_size?: number
}

export function getAuditLogs(params: AuditLogParams): Promise<PageResponse<AuditLog>> {
  return request.get('/audit-logs', { params })
}

export function getAuditLog(id: string): Promise<AuditLog> {
  return request.get(`/audit-logs/${id}`)
}