import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import WorkbenchAuthScreen from "@/pages/workbench/components/WorkbenchAuthScreen/WorkbenchAuthScreen.tsx"
import {ChatItem} from "@/pages/workbench/processes/chatReducer.ts"

describe("WorkbenchAuthScreen", () => {
    it("renders a sign-in CTA when an auth_link item is present", () => {
        const items: ChatItem[] = [{kind: "auth_link", key: "auth_link-0", url: "https://claude.ai/auth/abc"}]
        render(<WorkbenchAuthScreen items={items} status="open" onSubmitCode={vi.fn()}/>)

        const link = screen.getByRole("link", {name: "Sign in with Claude"})
        expect(link).toHaveAttribute("href", "https://claude.ai/auth/abc")
    })

    it("shows a waiting state when no auth_link item is present yet", () => {
        render(<WorkbenchAuthScreen items={[]} status="open" onSubmitCode={vi.fn()}/>)

        expect(screen.getByText("Waiting for sign-in link…")).toBeInTheDocument()
        expect(screen.queryByRole("link", {name: "Sign in with Claude"})).not.toBeInTheDocument()
    })

    it("calls onSubmitCode with the typed code while a code prompt is pending", () => {
        const items: ChatItem[] = [{kind: "auth_code_needed", key: "auth_code_needed-0", resolved: false}]
        const onSubmitCode = vi.fn()
        render(<WorkbenchAuthScreen items={items} status="open" onSubmitCode={onSubmitCode}/>)

        const input = screen.getByPlaceholderText("Enter the authentication code…")
        fireEvent.change(input, {target: {value: "123456"}})
        fireEvent.click(screen.getByRole("button", {name: "Submit"}))

        expect(onSubmitCode).toHaveBeenCalledWith("123456")
    })

    it("disables the code form while a pending code prompt exists but the connection is not open", () => {
        const items: ChatItem[] = [{kind: "auth_code_needed", key: "auth_code_needed-0", resolved: false}]
        const onSubmitCode = vi.fn()
        render(<WorkbenchAuthScreen items={items} status="reconnecting" onSubmitCode={onSubmitCode}/>)

        // Type a code so the only remaining reason the button could be disabled is the
        // connection status, not an empty input.
        const input = screen.getByPlaceholderText("Enter the authentication code…")
        fireEvent.change(input, {target: {value: "123456"}})

        expect(screen.getByRole("button", {name: "Submit"})).toBeDisabled()

        fireEvent.click(screen.getByRole("button", {name: "Submit"}))
        expect(onSubmitCode).not.toHaveBeenCalled()
    })

    it("uses the latest auth_link when a stale one lingers from a prior container run", () => {
        const items: ChatItem[] = [
            {kind: "auth_link", key: "auth_link-0", url: "https://claude.ai/auth/stale"},
            {kind: "auth_link", key: "auth_link-1", url: "https://claude.ai/auth/current"},
        ]
        render(<WorkbenchAuthScreen items={items} status="open" onSubmitCode={vi.fn()}/>)

        const link = screen.getByRole("link", {name: "Sign in with Claude"})
        expect(link).toHaveAttribute("href", "https://claude.ai/auth/current")
    })

    it("surfaces the latest error item's text", () => {
        const items: ChatItem[] = [
            {kind: "error", key: "error-0", text: "first failure"},
            {kind: "error", key: "error-1", text: "second failure"},
        ]
        render(<WorkbenchAuthScreen items={items} status="open" onSubmitCode={vi.fn()}/>)

        expect(screen.getByText("second failure")).toBeInTheDocument()
        expect(screen.queryByText("first failure")).not.toBeInTheDocument()
    })
})
