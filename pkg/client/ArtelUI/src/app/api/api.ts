export const GrpcAuthHeader = "Grpc-Metadata-Authorization"

export interface InitReq {
    pathPrefix: string,
    headers: {
        "Grpc-Metadata-Authorization": string,
    },
}

export type Options = {
    accessToken: string
}

export function apiPrefix(opts?: Options): InitReq {
    return {
        pathPrefix: import.meta.env.VITE_ARTEL_API,
        headers: {
            "Grpc-Metadata-Authorization": opts ? opts.accessToken : "",
        }
    }
}
