export enum ErrorReason {
    ACCESS_TOKEN_NOT_FOUND = "ACCESS_TOKEN_NOT_FOUND",
    REFRESH_TOKEN_NOT_FOUND = "REFRESH_TOKEN_NOT_FOUND",
}

export enum Errors {
    PERMISSION_DENIED = 7,
    UNAUTHENTICATED = 16,
}

export interface GrpcErrorDetail {
    reason?: ErrorReason
}

export interface GrpcError {
    code: number
    message: string
    details: GrpcErrorDetail[]
}

export function isReason(details: GrpcErrorDetail[], reason: ErrorReason): boolean {
    return details.some(d => d.reason === reason)
}

type ServiceErrorOption = (e: ServiceError) => void

export function WithTitle(title: string): ServiceErrorOption {
    return (e: ServiceError) => { e.title = title }
}

export function WithIsNonRetryable(v: boolean): ServiceErrorOption {
    return (e: ServiceError) => { e.isNonRetryable = v }
}

export function WithReason(reason?: ErrorReason): ServiceErrorOption {
    return (e: ServiceError) => { e.reason = reason }
}

export class ServiceError extends Error {
    title: string = ""
    isNonRetryable: boolean = false
    reason?: ErrorReason

    constructor(...opts: ServiceErrorOption[]) {
        super()
        opts.forEach(o => o(this))
        this.message = this.title
    }
}
