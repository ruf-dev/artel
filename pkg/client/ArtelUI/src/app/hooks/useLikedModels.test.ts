import {createElement} from "react"
import type {ReactNode} from "react"
import {QueryClient, QueryClientProvider} from "@tanstack/react-query"
import {act, renderHook, waitFor} from "@testing-library/react"
import {beforeEach, describe, expect, it, vi} from "vitest"

const h = vi.hoisted(() => ({
    getUserSettings: vi.fn(),
    setLikedModels: vi.fn(),
    authenticated: true,
}))
const mockBakeError = vi.fn()

vi.mock("@/processes/UserSettings.ts", async importOriginal => {
    const actual = await importOriginal<typeof import("@/processes/UserSettings.ts")>()
    return {
        ...actual,
        userSettingsService: {
            getUserSettings: h.getUserSettings,
            setLikedModels: h.setLikedModels,
        },
    }
})

vi.mock("@/hooks/user/User.ts", () => ({
    default: () => ({auth: {isAuthenticated: () => h.authenticated}}),
}))

vi.mock("@/app/hooks/useErrorToast.ts", () => ({
    useBakeError: () => mockBakeError,
}))

import {useLikedModels} from "@/app/hooks/useLikedModels.ts"
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
    h.setLikedModels.mockResolvedValue(undefined)
})

describe("useLikedModels", () => {
    it("defaults to an empty liked list before the settings load", () => {
        const client = new QueryClient({defaultOptions: {queries: {retry: false}}})
        const {result} = renderHook(() => useLikedModels(), {wrapper: makeWrapper(client)})

        expect(result.current.likedIds).toEqual([])
        expect(result.current.isLiked("openai/gpt-4")).toBe(false)
    })

    it("does not fetch when the user is not authenticated", () => {
        h.authenticated = false
        const client = new QueryClient({defaultOptions: {queries: {retry: false}}})

        renderHook(() => useLikedModels(), {wrapper: makeWrapper(client)})

        expect(h.getUserSettings).not.toHaveBeenCalled()
    })

    it("reads the liked list from GetUserSettings", async () => {
        h.getUserSettings.mockResolvedValue({
            userPrompt: "",
            likedOpenrouterModels: ["openai/gpt-4"],
            lastUsedModel: undefined,
        })
        const client = new QueryClient({defaultOptions: {queries: {retry: false}}})

        const {result} = renderHook(() => useLikedModels(), {wrapper: makeWrapper(client)})

        await waitFor(() => expect(result.current.likedIds).toEqual(["openai/gpt-4"]))
        expect(result.current.isLiked("openai/gpt-4")).toBe(true)
        expect(result.current.isLiked("anthropic/claude")).toBe(false)
    })

    it("toggleLiked adds an id optimistically and calls SetLikedModels with the full array", async () => {
        const client = new QueryClient({defaultOptions: {queries: {retry: false}}})
        const {result} = renderHook(() => useLikedModels(), {wrapper: makeWrapper(client)})
        await waitFor(() => expect(client.getQueryData(userSettingsQueryKey())).toBeDefined())

        act(() => {
            result.current.toggleLiked("openai/gpt-4")
        })

        await waitFor(() => expect(result.current.likedIds).toEqual(["openai/gpt-4"]))
        await waitFor(() => expect(h.setLikedModels).toHaveBeenCalledWith(["openai/gpt-4"]))
    })

    it("toggleLiked removes an already-liked id", async () => {
        h.getUserSettings.mockResolvedValue({
            userPrompt: "",
            likedOpenrouterModels: ["openai/gpt-4", "anthropic/claude"],
            lastUsedModel: undefined,
        })
        const client = new QueryClient({defaultOptions: {queries: {retry: false}}})
        const {result} = renderHook(() => useLikedModels(), {wrapper: makeWrapper(client)})
        await waitFor(() => expect(result.current.likedIds).toEqual(["openai/gpt-4", "anthropic/claude"]))

        act(() => {
            result.current.toggleLiked("openai/gpt-4")
        })

        await waitFor(() => expect(result.current.likedIds).toEqual(["anthropic/claude"]))
        await waitFor(() => expect(h.setLikedModels).toHaveBeenCalledWith(["anthropic/claude"]))
    })

    it("rolls back the optimistic update and bakes an error when SetLikedModels fails", async () => {
        h.setLikedModels.mockRejectedValue(new Error("boom"))
        const client = new QueryClient({defaultOptions: {queries: {retry: false}}})
        const {result} = renderHook(() => useLikedModels(), {wrapper: makeWrapper(client)})
        await waitFor(() => expect(client.getQueryData(userSettingsQueryKey())).toBeDefined())

        act(() => {
            result.current.toggleLiked("openai/gpt-4")
        })

        await waitFor(() =>
            expect(mockBakeError).toHaveBeenCalledWith("Failed to update liked models", expect.any(Error)))
        expect(result.current.likedIds).toEqual([])
    })
})
