import {render, screen, waitFor} from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import {afterEach, beforeEach, describe, expect, it, vi} from "vitest"

import ChatHistoryDialog from "@/dialogs/ChatHistoryDialog/ChatHistoryDialog.tsx"
import * as chatHistoryApi from "@/pages/workbench/processes/chatHistoryApi.ts"
import {ChatEvent} from "@/pages/workbench/processes/chatProtocol.ts"

describe("ChatHistoryDialog", () => {
    let listSessionsSpy: ReturnType<typeof vi.spyOn>
    let getSessionSpy: ReturnType<typeof vi.spyOn>

    beforeEach(() => {
        listSessionsSpy = vi.spyOn(chatHistoryApi, "listChatSessions").mockResolvedValue([])
        getSessionSpy = vi.spyOn(chatHistoryApi, "getChatSession").mockResolvedValue([])
    })

    afterEach(() => {
        vi.restoreAllMocks()
    })

    it("displays 'no sessions' when list is empty", async () => {
        render(<ChatHistoryDialog vaultId="v1"/>)

        await waitFor(() => {
            expect(screen.getByText("No previous chat sessions yet.")).toBeInTheDocument()
        })
    })

    it("fetches and displays session list on mount", async () => {
        const sessions = [
            {id: "s1", firstUserMessage: "Hello", lastActivityAt: "2026-01-01T12:00:00Z"},
            {id: "s2", firstUserMessage: "Another chat", lastActivityAt: "2026-01-02T12:00:00Z"},
        ]
        listSessionsSpy.mockResolvedValueOnce(sessions)

        render(<ChatHistoryDialog vaultId="v1"/>)

        await waitFor(() => {
            expect(screen.getByText("Hello")).toBeInTheDocument()
            expect(screen.getByText("Another chat")).toBeInTheDocument()
        })
    })

    it("displays '(empty session)' when a session has no firstUserMessage", async () => {
        const sessions = [
            {id: "s1", lastActivityAt: "2026-01-01T12:00:00Z"},
        ]
        listSessionsSpy.mockResolvedValueOnce(sessions)

        render(<ChatHistoryDialog vaultId="v1"/>)

        await waitFor(() => {
            expect(screen.getByText("(empty session)")).toBeInTheDocument()
        })
    })

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

        render(<ChatHistoryDialog vaultId="v1"/>)

        const button = await waitFor(() => screen.getByText("Test"))
        await userEvent.click(button)

        // Should load the session and display the message content
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

        render(<ChatHistoryDialog vaultId="v1"/>)

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
})
