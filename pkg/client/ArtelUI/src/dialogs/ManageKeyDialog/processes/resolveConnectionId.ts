import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"

export const ARTEL_BOT_OPTION = "__artel_bot__"

// Turns the SelectConnectionScreen selection into a real external-connection id:
// the "Artel bot" sentinel materializes a managed telegram connection on demand,
// every other value passes through unchanged.
export async function resolveConnectionId(
    selectedId: string,
    onError: (err: unknown) => void,
): Promise<string> {
    if (selectedId !== ARTEL_BOT_OPTION) {
        return selectedId
    }
    return useExternalConnections.getState()
        .ensureArtelTelegramConnection()
        .catch(err => {
            onError(err)
            return ""
        })
}
