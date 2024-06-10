/// <reference types="vite/client" />

export interface Frontend {
  id: number
  created_at: string
  updated_at: string
  interface: string
  port: number
  domain: string
  backends: Backend[]
  aliases: Alias[]
}

export interface Backend {
  id: number
  created_at: string
  updated_at: string
  address: string
  https: boolean
  https_verify: boolean
}

export interface Alias {
  id: number
  created_at: string
  updated_at: string
  domain: string
}


export interface FrontendStatus {
  frontend_id: number
  Servers: BackendStatus[]
  bytes_in_total: number
  bytes_out_total: number
  requests_total: number
  responses_total_1xx: number
  responses_total_2xx: number
  responses_total_3xx: number
  responses_total_4xx: number
  responses_total_5xx: number
  responses_total_other: number
}

export interface BackendStatus {
  server_id: number
  address: string
  sockerr: number
  l4ok: number
  l4tout: number
  l4con: number
  l6ok: number
  l6tout: number
  l6rsp: number
  l7tout: number
  l7rsp: number
  l7ok: number
  l7okc: number
  l7sts: number
}


export interface HaproxyLog {
  id: number
  created_at: string
  updated_at: string
  data: string
}

export interface Certificate {
  id: number
  created_at: string
  updated_at: string
  name: string
  pem_path: string
}