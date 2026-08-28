import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import WorkbenchTopbar from "@/pages/workbench/components/WorkbenchTopbar/WorkbenchTopbar.tsx"

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

function makeCtx(tweaksOpen = false) {
    return {tweaksOpen, tweaksSection: undefined, openTweaks: vi.fn(), closeTweaks: vi.fn()}
}

const base = {
    effectiveMode: "api" as const,
    exists: false,
    status: "running",
    vaultId: "v1",
    view: "chat" as const,
    onViewChange: vi.fn(),
    chatLocked: false,
    navOpen: true,
    onToggleNav: vi.fn(),
    onStart: vi.fn(),
    onStop: vi.fn(),
    stopping: false,
    starting: false,
    model,
    ctx: makeCtx(),
}

describe("WorkbenchTopbar", () => {
    it("api mode: live model switcher, no badge / start-stop / view switch, live tweaks toggle", () => {
        render(<WorkbenchTopbar {...base} effectiveMode="api"/>)

        expect(screen.getByTestId("model-switcher")).toHaveAttribute("data-disabled", "false")
        expect(screen.queryByText("Running")).not.toBeInTheDocument()
        expect(screen.queryByRole("button", {name: "Stop"})).not.toBeInTheDocument()
        expect(screen.queryByRole("button", {name: "Chat"})).not.toBeInTheDocument()
        expect(screen.getByRole("button", {name: "Tweaks"})).not.toHaveAttribute("aria-disabled")
    })

    it("clicking Tweaks opens the panel when it's closed", () => {
        const ctx = makeCtx(false)
        render(<WorkbenchTopbar {...base} effectiveMode="api" ctx={ctx}/>)

        fireEvent.click(screen.getByRole("button", {name: "Tweaks"}))

        expect(ctx.openTweaks).toHaveBeenCalledTimes(1)
        expect(ctx.closeTweaks).not.toHaveBeenCalled()
    })

    it("clicking Tweaks closes the panel when it's open", () => {
        const ctx = makeCtx(true)
        render(<WorkbenchTopbar {...base} effectiveMode="api" ctx={ctx}/>)

        fireEvent.click(screen.getByRole("button", {name: "Tweaks"}))

        expect(ctx.closeTweaks).toHaveBeenCalledTimes(1)
        expect(ctx.openTweaks).not.toHaveBeenCalled()
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

    it("fires onToggleNav when the nav button is clicked", () => {
        const onToggleNav = vi.fn()
        render(<WorkbenchTopbar {...base} effectiveMode="api" onToggleNav={onToggleNav}/>)

        fireEvent.click(screen.getByRole("button", {name: "Toggle conversations"}))

        expect(onToggleNav).toHaveBeenCalledTimes(1)
    })

    it("marks the Chat segment aria-disabled when chatLocked", () => {
        render(<WorkbenchTopbar {...base} effectiveMode="docker" exists status="running" chatLocked/>)

        expect(screen.getByRole("button", {name: "Chat"})).toHaveAttribute("aria-disabled", "true")
    })
})
