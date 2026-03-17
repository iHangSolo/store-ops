import request from '@/utils/request'
import type { PageResponse } from '@/utils/request'

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResult {
  token: string
  expires_at: string
  admin: {
    id: string
    username: string
    nickname: string
    role: string
  }
}

export function login(data: LoginRequest): Promise<LoginResult> {
  return request.post('/auth/login', data)
}

export function logout(): Promise<void> {
  return request.post('/auth/logout')
}