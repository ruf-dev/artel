import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import TopbarRight from "@/pages/workbench/components/WorkbenchTopbar/components/TopbarRight/TopbarRight.tsx"

vi.mock("morphicons/react", () => ({MorphIcon: () => <span data-testid="icon"/>}))

vi.mock(
    "@/pages/workbench/components/WorkbenchTopbar/components/StartStopButton/StartStopButton.tsx",
    () => ({default: () => <span data-testid="start-stop-button"/>}),
)

vi.mock(
    "@/pages/workbench/components/WorkbenchTopbar/components/WorkbenchSettingsMenu/WorkbenchSettingsMenu.tsx",
    () => ({default: () => <span data-testid="settings-menu"/>}),
)

vi.mock(
    "@/components/atoms/SegmentedControl/SegmentedControl.tsx",
    () => ({default: () => <span data-testid="segmented-control"/>}),
)

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
    onNewChat: vi.fn(),
}

describe("TopbarRight", () => {
    it("always renders the New Chat button", () => {
        render(<TopbarRight {...base} effectiveMode="api"/>)

        const newChatBtn = screen.getByRole("button", {name: "New chat"})
        expect(newChatBtn).toBeInTheDocument()
    })

    it("New Chat button calls onNewChat when clicked", () => {
        const onNewChat = vi.fn()
        render(<TopbarRight {...base} onNewChat={onNewChat}/>)

        const newChatBtn = screen.getByRole("button", {name: "New chat"})
        fireEvent.click(newChatBtn)

        expect(onNewChat).toHaveBeenCalledTimes(1)
    })

    it("api mode: no view switch, no start-stop, no settings", () => {
        render(<TopbarRight {...base} effectiveMode="api" exists={false}/>)

        expect(screen.queryByTestId("segmented-control")).not.toBeInTheDocument()
        expect(screen.queryByTestId("start-stop-button")).not.toBeInTheDocument()
        expect(screen.queryByTestId("settings-menu")).not.toBeInTheDocument()
    })

    it("docker running + exists: view switch, start-stop, settings", () => {
        render(<TopbarRight {...base} effectiveMode="docker" exists status="running" vaultId="v1"/>)

        expect(screen.getByTestId("segmented-control")).toBeInTheDocument()
        expect(screen.getByTestId("start-stop-button")).toBeInTheDocument()
        expect(screen.getByTestId("settings-menu")).toBeInTheDocument()
    })

    it("docker exists but not running: no view switch, has start-stop, has settings", () => {
        render(<TopbarRight {...base} effectiveMode="docker" exists status="stopped" vaultId="v1"/>)

        expect(screen.queryByTestId("segmented-control")).not.toBeInTheDocument()
        expect(screen.getByTestId("start-stop-button")).toBeInTheDocument()
        expect(screen.getByTestId("settings-menu")).toBeInTheDocument()
    })

    it("docker without vaultId: no settings menu even if exists", () => {
        render(<TopbarRight {...base} effectiveMode="docker" exists status="running" vaultId={undefined}/>)

        expect(screen.queryByTestId("settings-menu")).not.toBeInTheDocument()
    })
})
