export const GrpcAuthHeader = "Grpc-Metadata-Authorization"

export interface TelegramLoginResponse {
    token: string
    expiresAt: string
}


export interface InitReq {
    pathPrefix: string,
    headers: {
        "Content-Type"?: string,
        "Grpc-Metadata-Authorization": string,
    },
}

export type Options = {
    accessToken: string
}

export function apiPrefix(opts?: Options): InitReq {
    return {
        pathPrefix: "",
        headers: {
            "Content-Type": "application/json",
            "Grpc-Metadata-Authorization": opts ? opts.accessToken : "",
        }
    }
}
