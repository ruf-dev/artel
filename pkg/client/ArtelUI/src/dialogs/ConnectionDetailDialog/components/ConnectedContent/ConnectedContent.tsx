import {ExternalConnectionInfo, ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import GenericConnectionContent from "@/widgets/GenericConnectionContent/GenericConnectionContent.tsx"
import GoogleConnectionContent from "@/widgets/GoogleConnectionContent/GoogleConnectionContent.tsx"
import GoogleSheetsConnectionContent from "@/widgets/GoogleSheetsConnectionContent/GoogleSheetsConnectionContent.tsx"

interface ConnectedContentProps {
    provider: ExternalProvider
    connection: ExternalConnectionInfo
    onDisconnect: () => void
}

export default function ConnectedContent({provider, connection, onDisconnect}: ConnectedContentProps) {
    if (connection.google && provider === ExternalProvider.EXTERNAL_PROVIDER_GOOGLE_SHEETS) {
        return <GoogleSheetsConnectionContent connection={connection} onDisconnect={onDisconnect}/>
    }
    if (connection.google) {
        return <GoogleConnectionContent connection={connection} onDisconnect={onDisconnect}/>
    }
    return <GenericConnectionContent connection={connection} onDisconnect={onDisconnect}/>
}
