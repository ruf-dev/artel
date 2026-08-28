import {createElement} from "react"
import type {ReactNode} from "react"
import {QueryClient, QueryClientProvider} from "@tanstack/react-query"
import {renderHook, waitFor} from "@testing-library/react"
import {beforeEach, describe, expect, it, vi} from "vitest"

const h = vi.hoisted(() => ({
    listFolders: vi.fn(),
    listNotes: vi.fn(),
    authenticated: true,
}))

vi.mock("@/processes/Notes.ts", () => ({
    notesService: {
        listFolders: h.listFolders,
        listNotes: h.listNotes,
    },
}))

vi.mock("@/hooks/user/User.ts", () => ({
    default: () => ({auth: {isAuthenticated: () => h.authenticated}}),
}))

import {useWorkbenchVaultFiles} from "@/pages/workbench/processes/useWorkbenchVaultFiles.ts"

function wrapper({children}: {children: ReactNode}) {
    const client = new QueryClient({defaultOptions: {queries: {retry: false}}})
    return createElement(QueryClientProvider, {client}, children)
}

beforeEach(() => {
    vi.clearAllMocks()
    h.authenticated = true
    h.listFolders.mockResolvedValue(["notes", "notes/sub"])
    h.listNotes.mockResolvedValue([{path: "notes/a.md"}, {path: "b.md"}])
})

describe("useWorkbenchVaultFiles", () => {
    it("does not fetch while vaultId is undefined", () => {
        renderHook(() => useWorkbenchVaultFiles(undefined), {wrapper})

        expect(h.listFolders).not.toHaveBeenCalled()
        expect(h.listNotes).not.toHaveBeenCalled()
    })

    it("does not fetch when the user is not authenticated", () => {
        h.authenticated = false

        renderHook(() => useWorkbenchVaultFiles("v1"), {wrapper})

        expect(h.listFolders).not.toHaveBeenCalled()
        expect(h.listNotes).not.toHaveBeenCalled()
    })

    it("loads folders and notes for the vault", async () => {
        const {result} = renderHook(() => useWorkbenchVaultFiles("v1"), {wrapper})

        await waitFor(() => expect(result.current.isLoading).toBe(false))

        expect(h.listFolders).toHaveBeenCalledWith("v1")
        expect(h.listNotes).toHaveBeenCalledWith("v1")
        expect(result.current.folders).toEqual(["notes", "notes/sub"])
        expect(result.current.notes).toEqual([{path: "notes/a.md"}, {path: "b.md"}])
        expect(result.current.error).toBeFalsy()
    })

    it("surfaces the first query error", async () => {
        const boom = new Error("boom")
        h.listFolders.mockRejectedValue(boom)

        const {result} = renderHook(() => useWorkbenchVaultFiles("v1"), {wrapper})

        await waitFor(() => expect(result.current.error).toBe(boom))
    })
})
