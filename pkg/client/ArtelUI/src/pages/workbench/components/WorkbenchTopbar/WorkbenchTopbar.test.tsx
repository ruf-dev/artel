import {describe, expect, it, vi} from "vitest"
import {render, screen} from "@testing-library/react"

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
})
