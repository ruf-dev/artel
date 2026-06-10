export const DetailType = {
    PreconditionFailure: "type.googleapis.com/google.rpc.PreconditionFailure",
    BadRequest:          "type.googleapis.com/google.rpc.BadRequest",
    ErrorInfo:           "type.googleapis.com/google.rpc.ErrorInfo",
    RetryInfo:           "type.googleapis.com/google.rpc.RetryInfo",
    ResourceInfo:        "type.googleapis.com/google.rpc.ResourceInfo",
    Help:                "type.googleapis.com/google.rpc.Help",
    LocalizedMessage:    "type.googleapis.com/google.rpc.LocalizedMessage",
} as const

export type DetailTypeName = typeof DetailType[keyof typeof DetailType]

export interface PreconditionViolation {
    type: string
    subject: string
    description: string
}

export interface FieldViolation {
    field: string
    description: string
}

export interface PreconditionFailureDetail {
    "@type": typeof DetailType.PreconditionFailure
    violations: PreconditionViolation[]
}

export interface BadRequestDetail {
    "@type": typeof DetailType.BadRequest
    fieldViolations: FieldViolation[]
}

export interface ErrorInfoDetail {
    "@type": typeof DetailType.ErrorInfo
    reason: string
    domain: string
    metadata: Record<string, string>
}

export interface RetryInfoDetail {
    "@type": typeof DetailType.RetryInfo
    retryDelay: { seconds: string; nanos: number }
}

export interface ResourceInfoDetail {
    "@type": typeof DetailType.ResourceInfo
    resourceType: string
    resourceName: string
    owner: string
    description: string
}

export interface HelpLink {
    description: string
    url: string
}

export interface HelpDetail {
    "@type": typeof DetailType.Help
    links: HelpLink[]
}

export interface LocalizedMessageDetail {
    "@type": typeof DetailType.LocalizedMessage
    locale: string
    message: string
}

export type GrpcErrorDetail =
    | PreconditionFailureDetail
    | BadRequestDetail
    | ErrorInfoDetail
    | RetryInfoDetail
    | ResourceInfoDetail
    | HelpDetail
    | LocalizedMessageDetail
    | { "@type": string; [key: string]: unknown }

export interface GrpcStatusError {
    code: number
    message: string
    details: GrpcErrorDetail[]
}

export function isGrpcError(err: unknown): err is GrpcStatusError {
    return (
        typeof err === "object" &&
        err !== null &&
        "code" in err &&
        typeof (err as GrpcStatusError).code === "number" &&
        "message" in err &&
        typeof (err as GrpcStatusError).message === "string" &&
        "details" in err &&
        Array.isArray((err as GrpcStatusError).details)
    )
}

export function getDetail<T extends GrpcErrorDetail>(
    err: GrpcStatusError,
    type: DetailTypeName,
): T | undefined {
    return err.details.find(d => d["@type"] === type) as T | undefined
}

export function getPreconditionViolations(err: GrpcStatusError): PreconditionViolation[] {
    const d = getDetail<PreconditionFailureDetail>(err, DetailType.PreconditionFailure)
    return d?.violations ?? []
}

export function getFieldViolations(err: GrpcStatusError): FieldViolation[] {
    const d = getDetail<BadRequestDetail>(err, DetailType.BadRequest)
    return d?.fieldViolations ?? []
}

export function grpcErrorMessage(err: GrpcStatusError): string {
    const localized = getDetail<LocalizedMessageDetail>(err, DetailType.LocalizedMessage)
    if (localized) return localized.message
    return err.message
}

export function retryOnStatus(...codes: number[]) {
    return (_: number, err: unknown): boolean =>
        isGrpcError(err) && codes.includes(err.code)
}

export function notRetryOnStatus(...codes: number[]) {
    return (_: number, err: unknown): boolean =>
        !isGrpcError(err) || !codes.includes(err.code)
}
