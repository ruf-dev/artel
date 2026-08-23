import {render, screen, waitFor} from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import {afterEach, beforeEach, describe, expect, it, vi} from "vitest"

import ChatHistorySidebar from "@/pages/workbench/components/Chat/components/ChatHistorySidebar/ChatHistorySidebar.tsx"
import * as chatHistoryApi from "@/pages/workbench/processes/chatHistoryApi.ts"
import {ChatEvent} from "@/pages/workbench/processes/chatProtocol.ts"

let listSessionsSpy: ReturnType<typeof vi.spyOn>
let getSessionSpy: ReturnType<typeof vi.spyOn>

beforeEach(() => {
    listSessionsSpy = vi.spyOn(chatHistoryApi, "listChatSessions").mockResolvedValue([])
    getSessionSpy = vi.spyOn(chatHistoryApi, "getChatSession").mockResolvedValue([])
})

afterEach(() => {
    vi.restoreAllMocks()
})

describe("ChatHistorySidebar - open/closed fetching", () => {
    it("does not fetch sessions while closed", () => {
        render(<ChatHistorySidebar vaultId="v1" open={false} onClose={vi.fn()}/>)

        expect(listSessionsSpy).not.toHaveBeenCalled()
    })

    it("fetches sessions once opened", async () => {
        const sessions = [
            {id: "s1", firstUserMessage: "Hello", lastActivityAt: "2026-01-01T12:00:00Z"},
        ]
        listSessionsSpy.mockResolvedValueOnce(sessions)

        const {rerender} = render(<ChatHistorySidebar vaultId="v1" open={false} onClose={vi.fn()}/>)
        expect(listSessionsSpy).not.toHaveBeenCalled()

        rerender(<ChatHistorySidebar vaultId="v1" open onClose={vi.fn()}/>)

        await waitFor(() => {
            expect(listSessionsSpy).toHaveBeenCalledWith("v1")
            expect(screen.getByText("Hello")).toBeInTheDocument()
        })
    })

    it("calls onClose when the close button is clicked", async () => {
        const onClose = vi.fn()
        render(<ChatHistorySidebar vaultId="v1" open onClose={onClose}/>)

        await waitFor(() => {
            expect(screen.getByText("No previous chat sessions yet.")).toBeInTheDocument()
        })

        await userEvent.click(screen.getByLabelText("Close chat history"))

        expect(onClose).toHaveBeenCalledTimes(1)
    })
})

describe("ChatHistorySidebar - list screen", () => {
    it("displays 'no sessions' when list is empty", async () => {
        render(<ChatHistorySidebar vaultId="v1" open onClose={vi.fn()}/>)

        await waitFor(() => {
            expect(screen.getByText("No previous chat sessions yet.")).toBeInTheDocument()
        })
    })

    it("displays '(empty session)' when a session has no firstUserMessage", async () => {
        const sessions = [
            {id: "s1", lastActivityAt: "2026-01-01T12:00:00Z"},
        ]
        listSessionsSpy.mockResolvedValueOnce(sessions)

        render(<ChatHistorySidebar vaultId="v1" open onClose={vi.fn()}/>)

        await waitFor(() => {
            expect(screen.getByText("(empty session)")).toBeInTheDocument()
        })
    })
})

describe("ChatHistorySidebar - detail screen", () => {
    it("loads and displays session detail when a session is clicked", async () => {
        const sessions = [
            {id: "s1", firstUserMessage: "Test", lastActivityAt: "2026-01-01T12:00:00Z"},
        ]
        const sessionEvents: ChatEvent[] = [
            {type: "user_message", text: "Test", seq: 1},
            {type: "assistant_text_done", id: "a1", text: "Response", seq: 2},
        ]

        listSessionsSpy.mockResolvedValueOnce(sessions)
        getSessionSpy.mockResolvedValueOnce(sessionEvents)

        render(<ChatHistorySidebar vaultId="v1" open onClose={vi.fn()}/>)

        const button = await waitFor(() => screen.getByText("Test"))
        await userEvent.click(button)

        await waitFor(() => {
            expect(screen.getByText("Response")).toBeInTheDocument()
        })
    })

    it("returns to list when back button is clicked", async () => {
        const sessions = [
            {id: "s1", firstUserMessage: "Test", lastActivityAt: "2026-01-01T12:00:00Z"},
        ]
        const sessionEvents: ChatEvent[] = [
            {type: "assistant_text_done", id: "a1", text: "Response", seq: 1},
        ]

        listSessionsSpy.mockResolvedValueOnce(sessions)
        getSessionSpy.mockResolvedValueOnce(sessionEvents)

        render(<ChatHistorySidebar vaultId="v1" open onClose={vi.fn()}/>)

        const button = await waitFor(() => screen.getByText("Test"))
        await userEvent.click(button)

        await waitFor(() => {
            expect(screen.getByText("Response")).toBeInTheDocument()
        })

        const backButton = screen.getByRole("button", {name: "Back to list"})
        await userEvent.click(backButton)

        await waitFor(() => {
            expect(screen.getByText("Chat History")).toBeInTheDocument()
        })
    })

    it("resets back to the list screen the next time it is opened", async () => {
        const sessions = [
            {id: "s1", firstUserMessage: "Test", lastActivityAt: "2026-01-01T12:00:00Z"},
        ]
        const sessionEvents: ChatEvent[] = [
            {type: "assistant_text_done", id: "a1", text: "Response", seq: 1},
        ]

        listSessionsSpy.mockResolvedValue(sessions)
        getSessionSpy.mockResolvedValueOnce(sessionEvents)

        const {rerender} = render(<ChatHistorySidebar vaultId="v1" open onClose={vi.fn()}/>)

        const button = await waitFor(() => screen.getByText("Test"))
        await userEvent.click(button)

        await waitFor(() => {
            expect(screen.getByText("Response")).toBeInTheDocument()
        })

        // Close, then reopen — should be back on the list screen, not stuck on detail.
        rerender(<ChatHistorySidebar vaultId="v1" open={false} onClose={vi.fn()}/>)
        rerender(<ChatHistorySidebar vaultId="v1" open onClose={vi.fn()}/>)

        await waitFor(() => {
            expect(screen.getByText("Chat History")).toBeInTheDocument()
        })
    })
})
