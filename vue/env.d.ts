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
}

export interface Alias {
    id: number
    created_at: string
    updated_at: string
    domain: string
}