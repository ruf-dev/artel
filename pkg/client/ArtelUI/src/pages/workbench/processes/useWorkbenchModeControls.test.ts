import {renderHook, waitFor} from "@testing-library/react"
import {describe, expect, it, vi} from "vitest"

import {useWorkbenchModeControls} from "@/pages/workbench/processes/useWorkbenchModeControls.ts"
import type {SimpleChat} from "@/processes/SimpleChat.ts"

const chat = {id: "chat-1"} as SimpleChat

describe("useWorkbenchModeControls", () => {
    it("restores Simple Chat and its most recent thread from persisted chats", async () => {
        const {result} = renderHook(() => useWorkbenchModeControls({
            exists: false,
            handleCreateDocker: vi.fn(),
            simpleChats: [chat],
            simpleChatsLoading: false,
        }))

        await waitFor(() => expect(result.current.effectiveMode).toBe("api"))
        expect(result.current.simpleChatId).toBe("chat-1")
    })

    it("keeps the Docker fallback for an existing workbench without a Simple Chat", async () => {
        const {result} = renderHook(() => useWorkbenchModeControls({
            exists: true,
            handleCreateDocker: vi.fn(),
            simpleChats: [],
            simpleChatsLoading: false,
        }))

        await waitFor(() => expect(result.current.effectiveMode).toBe("docker"))
    })

    it("does not show the picker while persisted chats are loading", () => {
        const {result} = renderHook(() => useWorkbenchModeControls({
            exists: false,
            handleCreateDocker: vi.fn(),
            simpleChats: [],
            simpleChatsLoading: true,
        }))

        expect(result.current.effectiveMode).toBe("picking")
    })
})
