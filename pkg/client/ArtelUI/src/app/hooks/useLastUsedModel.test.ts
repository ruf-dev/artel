import {createElement} from "react"
import type {ReactNode} from "react"
import {QueryClient, QueryClientProvider} from "@tanstack/react-query"
import {act, renderHook, waitFor} from "@testing-library/react"
import {beforeEach, describe, expect, it, vi} from "vitest"

const h = vi.hoisted(() => ({
    getUserSettings: vi.fn(),
    setLastUsedModel: vi.fn(),
    authenticated: true,
}))
const mockBakeError = vi.fn()

vi.mock("@/processes/UserSettings.ts", async importOriginal => {
    const actual = await importOriginal<typeof import("@/processes/UserSettings.ts")>()
    return {
        ...actual,
        userSettingsService: {
            getUserSettings: h.getUserSettings,
            setLastUsedModel: h.setLastUsedModel,
        },
    }
})

vi.mock("@/hooks/user/User.ts", () => ({
    default: () => ({auth: {isAuthenticated: () => h.authenticated}}),
}))

vi.mock("@/app/hooks/useErrorToast.ts", () => ({
    useBakeError: () => mockBakeError,
}))

import {useLastUsedModel} from "@/app/hooks/useLastUsedModel.ts"
import {userSettingsQueryKey} from "@/processes/UserSettings.ts"

function makeWrapper(client: QueryClient) {
    return function Wrapper({children}: {children: ReactNode}) {
        return createElement(QueryClientProvider, {client}, children)
    }
}

beforeEach(() => {
    vi.clearAllMocks()
    h.authenticated = true
    h.getUserSettings.mockResolvedValue({userPrompt: "", likedOpenrouterModels: [], lastUsedModel: undefined})
    h.setLastUsedModel.mockResolvedValue(undefined)
})

describe("useLastUsedModel", () => {
    it("defaults to undefined before the settings load", () => {
        const client = new QueryClient({defaultOptions: {queries: {retry: false}}})
        const {result} = renderHook(() => useLastUsedModel(), {wrapper: makeWrapper(client)})

        expect(result.current.lastUsedModel).toBeUndefined()
    })

    it("does not fetch when the user is not authenticated", () => {
        h.authenticated = false
        const client = new QueryClient({defaultOptions: {queries: {retry: false}}})

        renderHook(() => useLastUsedModel(), {wrapper: makeWrapper(client)})

        expect(h.getUserSettings).not.toHaveBeenCalled()
    })

    it("reads the last used model from GetUserSettings", async () => {
        h.getUserSettings.mockResolvedValue({
            userPrompt: "",
            likedOpenrouterModels: [],
            lastUsedModel: "openai/gpt-4",
        })
        const client = new QueryClient({defaultOptions: {queries: {retry: false}}})

        const {result} = renderHook(() => useLastUsedModel(), {wrapper: makeWrapper(client)})

        await waitFor(() => expect(result.current.lastUsedModel).toBe("openai/gpt-4"))
    })

    it("setLastUsedModel updates the cache immediately and calls SetLastUsedModel", async () => {
        const client = new QueryClient({defaultOptions: {queries: {retry: false}}})
        const {result} = renderHook(() => useLastUsedModel(), {wrapper: makeWrapper(client)})
        await waitFor(() => expect(client.getQueryData(userSettingsQueryKey())).toBeDefined())

        act(() => {
            result.current.setLastUsedModel("anthropic/claude")
        })

        await waitFor(() => expect(result.current.lastUsedModel).toBe("anthropic/claude"))
        await waitFor(() => expect(h.setLastUsedModel).toHaveBeenCalledWith("anthropic/claude"))
    })

    it("bakes an error when SetLastUsedModel fails", async () => {
        h.setLastUsedModel.mockRejectedValue(new Error("boom"))
        const client = new QueryClient({defaultOptions: {queries: {retry: false}}})
        const {result} = renderHook(() => useLastUsedModel(), {wrapper: makeWrapper(client)})
        await waitFor(() => expect(client.getQueryData(userSettingsQueryKey())).toBeDefined())

        act(() => {
            result.current.setLastUsedModel("anthropic/claude")
        })

        await waitFor(() =>
            expect(mockBakeError).toHaveBeenCalledWith("Failed to save last used model", expect.any(Error)))
    })
})
