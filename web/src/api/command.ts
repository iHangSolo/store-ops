import request from '@/utils/request'

export interface CommandRequest {
  command: string
  timeout?: number
}

export interface CommandResult {
  id: string
  store_id: string
  command: string
  output: string
  status: 'pending' | 'running' | 'success' | 'failed' | 'timeout'
  duration_ms: number
  created_at: string
  completed_at: string | null
}

export function executeCommand(storeId: string, data: CommandRequest): Promise<{ command_id: string }> {
  return request.post(`/stores/${storeId}/command`, data)
}

export function getCommandResult(storeId: string, commandId: string): Promise<CommandResult> {
  return request.get(`/stores/${storeId}/command/${commandId}`)
}