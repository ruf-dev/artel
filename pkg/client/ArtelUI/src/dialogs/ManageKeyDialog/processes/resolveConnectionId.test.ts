import {describe, expect, it, vi} from "vitest"

import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {ARTEL_BOT_OPTION, resolveConnectionId} from "@/dialogs/ManageKeyDialog/processes/resolveConnectionId.ts"

vi.mock("@/app/hooks/ExternalConnections.ts", () => ({
    useExternalConnections: {getState: vi.fn()},
}))

function mockEnsure(impl: () => Promise<string>): void {
    vi.mocked(useExternalConnections.getState).mockReturnValue(
        {ensureArtelTelegramConnection: impl} as unknown as ReturnType<typeof useExternalConnections.getState>,
    )
}

describe("resolveConnectionId", () => {
    it("passes a real connection id straight through without touching the store", async () => {
        const onError = vi.fn()
        const got = await resolveConnectionId("conn-123", onError)
        expect(got).toBe("conn-123")
        expect(onError).not.toHaveBeenCalled()
    })

    it("materializes the managed connection when the Artel bot sentinel is selected", async () => {
        mockEnsure(() => Promise.resolve("managed-1"))
        const onError = vi.fn()
        const got = await resolveConnectionId(ARTEL_BOT_OPTION, onError)
        expect(got).toBe("managed-1")
        expect(onError).not.toHaveBeenCalled()
    })

    it("reports the error and returns an empty id when materializing fails", async () => {
        const err = new Error("no telegram identity")
        mockEnsure(() => Promise.reject(err))
        const onError = vi.fn()
        const got = await resolveConnectionId(ARTEL_BOT_OPTION, onError)
        expect(got).toBe("")
        expect(onError).toHaveBeenCalledWith(err)
    })
})
