import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import WorkbenchTopbar from "@/pages/workbench/components/WorkbenchTopbar/WorkbenchTopbar.tsx"

vi.mock("morphicons/react", () => ({MorphIcon: () => null}))

interface MockModelSwitcherProps {
    disabled?: boolean
    placeholder?: string
}

vi.mock("@/pages/workbench/components/ModelSwitcher/ModelSwitcher.tsx", () => ({
    default: ({disabled, placeholder}: MockModelSwitcherProps) => (
        <span
            data-testid="model-switcher"
            data-disabled={String(Boolean(disabled))}
            data-placeholder={placeholder ?? ""}
        />
    ),
}))

vi.mock(
    "@/pages/workbench/components/WorkbenchTopbar/components/WorkbenchSettingsMenu/WorkbenchSettingsMenu.tsx",
    () => ({default: () => <span data-testid="settings-menu"/>}),
)

const model = {models: ["a/b"], value: "a/b", isLoading: false, onChange: vi.fn()}

const base = {
    effectiveMode: "api" as const,
    exists: false,
    status: "running",
    vaultId: "v1",
    view: "chat" as const,
    onViewChange: vi.fn(),
    chatLocked: false,
    onStart: vi.fn(),
    onStop: vi.fn(),
    stopping: false,
    starting: false,
    model,
    onToggleNav: vi.fn(),
    showNavToggle: true,
    onNewChat: vi.fn(),
    newChatDisabled: false,
}

describe("WorkbenchTopbar", () => {
    it("api mode: live model switcher, no badge / start-stop / view switch", () => {
        render(<WorkbenchTopbar {...base} effectiveMode="api"/>)

        expect(screen.getByTestId("model-switcher")).toHaveAttribute("data-disabled", "false")
        expect(screen.queryByText("Running")).not.toBeInTheDocument()
        expect(screen.queryByRole("button", {name: "Stop"})).not.toBeInTheDocument()
        expect(screen.queryByRole("button", {name: "Chat"})).not.toBeInTheDocument()
    })

    it("docker running + exists: badge, disabled switcher, Chat/Terminal switch, start/stop, settings", () => {
        render(<WorkbenchTopbar {...base} effectiveMode="docker" exists status="running"/>)

        expect(screen.getByText("Running")).toBeInTheDocument()
        expect(screen.getByTestId("model-switcher")).toHaveAttribute("data-disabled", "true")
        expect(screen.getByRole("button", {name: "Chat"})).toBeInTheDocument()
        expect(screen.getByRole("button", {name: "Terminal"})).toBeInTheDocument()
        expect(screen.getByRole("button", {name: "Stop"})).toBeInTheDocument()
        expect(screen.getByTestId("settings-menu")).toBeInTheDocument()
    })

    it("renders nothing for docker mode without a provisioned workbench", () => {
        const {container} = render(<WorkbenchTopbar {...base} effectiveMode="docker" exists={false}/>)

        expect(container).toBeEmptyDOMElement()
    })

    it("marks the Chat segment aria-disabled when chatLocked", () => {
        render(<WorkbenchTopbar {...base} effectiveMode="docker" exists status="running" chatLocked/>)

        expect(screen.getByRole("button", {name: "Chat"})).toHaveAttribute("aria-disabled", "true")
    })

    it("showNavToggle: renders the burger and forwards clicks to onToggleNav", () => {
        const onToggleNav = vi.fn()
        render(<WorkbenchTopbar {...base} showNavToggle onToggleNav={onToggleNav}/>)

        // hidden: true — the burger's container is `display: none` above the mobile breakpoint.
        const burger = screen.getByRole("button", {name: "Open menu", hidden: true})
        expect(burger).toBeInTheDocument()

        fireEvent.click(burger)
        expect(onToggleNav).toHaveBeenCalledTimes(1)
    })

    it("showNavToggle false: no burger button", () => {
        render(<WorkbenchTopbar {...base} showNavToggle={false}/>)

        expect(screen.queryByRole("button", {name: "Open menu", hidden: true})).toBeNull()
    })

    it("api mode: New Chat button renders and calls onNewChat on click", () => {
        const onNewChat = vi.fn()
        render(<WorkbenchTopbar {...base} effectiveMode="api" onNewChat={onNewChat}/>)

        const newChatBtn = screen.getByRole("button", {name: "New chat"})
        expect(newChatBtn).toBeInTheDocument()

        fireEvent.click(newChatBtn)
        expect(onNewChat).toHaveBeenCalledTimes(1)
    })

    it("docker mode: New Chat button renders and calls onNewChat on click", () => {
        const onNewChat = vi.fn()
        render(<WorkbenchTopbar {...base} effectiveMode="docker" exists status="running" onNewChat={onNewChat}/>)

        const newChatBtn = screen.getByRole("button", {name: "New chat"})
        expect(newChatBtn).toBeInTheDocument()

        fireEvent.click(newChatBtn)
        expect(onNewChat).toHaveBeenCalledTimes(1)
    })

    it("newChatDisabled disables the New Chat button", () => {
        render(<WorkbenchTopbar {...base} effectiveMode="api" newChatDisabled/>)

        expect(screen.getByRole("button", {name: "New chat"})).toBeDisabled()
    })
})
