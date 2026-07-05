import {ExternalConnectionInfo, ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"

export function connectionLabel(c: ExternalConnectionInfo): string {
    if (c.google) return c.google.email ?? "Google account"
    if (c.provider === ExternalProvider.EXTERNAL_PROVIDER_EMAIL) {
        return c.generic?.fields?.username ?? "Email account"
    }
    return c.provider?.replace("EXTERNAL_PROVIDER_", "").toLowerCase() ?? "Connection"
}
