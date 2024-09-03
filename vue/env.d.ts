/// <reference types="vite/client" />

export interface Frontend {
  id: number
  created_at: string
  updated_at: string
  interface: string
  port: number
  domain: string
  https: boolean
  bw_limit: number
  bw_limit_unit: number
  bw_period: number
  bw_out_limit: number
  bw_out_limit_unit: number
  bw_out_period: number
  rate_limit: number
  rate_period: number
  hard_rate_limit: number
  hard_rate_period: number
  http_check: boolean
  http_check_method: string
  http_check_path: string
  http_check_expect_status: number
  http_check_interval: number
  http_check_fail_after: number
  http_check_recover_after: number
  request_body_limit: number
  request_body_limit_unit: number
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

export interface Settings {
  acme_email: string
  acme: SettingsAcme
}

export interface SettingsAcme {
  cf: SettingsAcmeCf[]
}

export interface SettingsAcmeCf {
  id: number
  name: string
  token: string
}

export interface HaproxyCrashReasonsData {
  has_crashed: boolean
  address_in_use: boolean
  address_in_use_log: string
  permission_denied_port: boolean
  permission_denied_port_log: string
}
