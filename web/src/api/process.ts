import request from '@/utils/request'

export interface Process {
  pid: number
  name: string
  cpu_percent: number
  memory_percent: number
  status: string
  create_time: number
}

export interface Service {
  name: string
  display_name: string
  status: 'running' | 'stopped' | 'paused'
  start_type: 'auto' | 'manual' | 'disabled'
}

export interface KillProcessRequest {
  pid: number
}

export function getProcesses(storeId: string): Promise<Process[]> {
  return request.get(`/stores/${storeId}/processes`)
}

export function killProcess(storeId: string, pid: number): Promise<void> {
  return request.delete(`/stores/${storeId}/processes`, { data: { pid } })
}

export function getServices(storeId: string): Promise<Service[]> {
  return request.get(`/stores/${storeId}/services`)
}

export function startService(storeId: string, name: string): Promise<void> {
  return request.post(`/stores/${storeId}/services/${name}/start`)
}

export function stopService(storeId: string, name: string): Promise<void> {
  return request.post(`/stores/${storeId}/services/${name}/stop`)
}