export const GrpcAuthHeader = "Grpc-Metadata-Authorization"

export interface TelegramLoginResponse {
    token: string
    expiresAt: string
}

export async function loginTelegram(idToken: string): Promise<TelegramLoginResponse> {
    const url = `${import.meta.env.VITE_ARTEL_API}/api/auth/login`
    const body = JSON.stringify({telegram: {id_token: idToken}})
    const resp = await fetch(url, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body,
    })
    if (!resp.ok) {
        const text = await resp.text()
        throw new Error(`telegram login failed: ${text}`)
    }
    return resp.json()
}

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
