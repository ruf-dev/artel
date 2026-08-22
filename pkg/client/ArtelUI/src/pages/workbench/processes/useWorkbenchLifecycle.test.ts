import {act, renderHook, waitFor} from "@testing-library/react"
import {beforeEach, describe, expect, it, vi} from "vitest"

// Mock localStorage before any module imports
const localStorageMock = {
    getItem: vi.fn(),
    setItem: vi.fn(),
    removeItem: vi.fn(),
    clear: vi.fn(),
}
Object.defineProperty(global, "localStorage", {
    value: localStorageMock,
    writable: true,
})

const mockStart = vi.fn()
const mockCreate = vi.fn(() => Promise.resolve())
const mockStop = vi.fn(() => Promise.resolve())
const mockBakeError = vi.fn()
const mockFetchConnections = vi.fn(() => Promise.resolve())
let mockConnections: {provider: string}[] = []

vi.mock("@/app/hooks/Workbench.ts", () => ({
    useWorkbenchMutations: () => ({
        start: mockStart,
        create: mockCreate,
        stop: mockStop,
    }),
}))

vi.mock("@/app/hooks/ExternalConnections.ts", () => ({
    useExternalConnections: () => ({
        connections: mockConnections,
        fetch: mockFetchConnections,
    }),
}))

vi.mock("@/app/hooks/useErrorToast.ts", () => ({
    useBakeError: () => mockBakeError,
}))

import {useWorkbenchLifecycle} from "@/pages/workbench/processes/useWorkbenchLifecycle.ts"
import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"

describe("useWorkbenchLifecycle", () => {
    beforeEach(() => {
        mockStart.mockClear()
        mockCreate.mockClear()
        mockStop.mockClear()
        mockBakeError.mockClear()
        mockFetchConnections.mockClear()
        mockConnections = []
        mockStart.mockReturnValue(Promise.resolve())
    })

    it("sets pendingAuthMode and startingWorkbench synchronously when handleStartWorkbench is called", async () => {
        let resolveStart!: () => void
        const startPromise = new Promise<void>((resolve) => {
            resolveStart = resolve
        })
        mockStart.mockReturnValue(startPromise)

        const {result} = renderHook(() => useWorkbenchLifecycle("vault-123"))

        expect(result.current.pendingAuthMode).toBeUndefined()
        expect(result.current.startingWorkbench).toBe(false)

        act(() => {
            void result.current.handleStartWorkbench("subscription_login")
        })

        // Synchronously, before the mutation resolves
        expect(result.current.pendingAuthMode).toBe("subscription_login")
        expect(result.current.startingWorkbench).toBe(true)

        act(() => {
            resolveStart()
        })

        await waitFor(() => expect(result.current.startingWorkbench).toBe(false))
    })

    it("clears pendingAuthMode and calls bakeError when start mutation rejects", async () => {
        const testError = new Error("start failed")
        let rejectStart!: (e: Error) => void
        const startPromise = new Promise<void>((_, reject) => {
            rejectStart = reject
        })
        mockStart.mockReturnValue(startPromise)

        const {result} = renderHook(() => useWorkbenchLifecycle("vault-123"))

        act(() => {
            void result.current.handleStartWorkbench("subscription_login")
        })

        expect(result.current.pendingAuthMode).toBe("subscription_login")

        act(() => {
            rejectStart(testError)
        })

        await waitFor(() => expect(result.current.pendingAuthMode).toBeUndefined())

        expect(mockBakeError).toHaveBeenCalledWith("Failed to start workbench", testError)
    })

    it("shows the auth-mode picker and does not auto-start when an Anthropic connection exists", async () => {
        mockConnections = [{provider: ExternalProvider.EXTERNAL_PROVIDER_ANTHROPIC}]

        const {result} = renderHook(() => useWorkbenchLifecycle("vault-123"))

        act(() => {
            result.current.handleCreate()
        })

        await waitFor(() => expect(result.current.showSetup).toBe(true))

        expect(mockStart).not.toHaveBeenCalled()
        expect(result.current.pendingAuthMode).toBeUndefined()
    })

    it("skips the picker and auto-starts subscription login when there is no Anthropic connection", async () => {
        mockConnections = []

        const {result} = renderHook(() => useWorkbenchLifecycle("vault-123"))

        act(() => {
            result.current.handleCreate()
        })

        await waitFor(() => expect(mockStart).toHaveBeenCalledWith("subscription_login"))

        expect(result.current.showSetup).toBe(false)
    })

    it("handleStartClick also skips the picker and auto-starts when there is no Anthropic connection", async () => {
        mockConnections = []

        const {result} = renderHook(() => useWorkbenchLifecycle("vault-123"))

        act(() => {
            result.current.handleStartClick()
        })

        await waitFor(() => expect(mockStart).toHaveBeenCalledWith("subscription_login"))

        expect(result.current.showSetup).toBe(false)
    })
})
