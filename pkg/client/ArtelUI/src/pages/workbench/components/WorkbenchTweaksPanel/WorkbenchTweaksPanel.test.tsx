import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import WorkbenchTweaksPanel from "@/pages/workbench/components/WorkbenchTweaksPanel/WorkbenchTweaksPanel.tsx"

vi.mock(
    // eslint-disable-next-line max-len -- path too long to wrap
    "@/pages/workbench/components/WorkbenchTweaksPanel/components/TweaksSystemPromptSection/TweaksSystemPromptSection.tsx",
    () => ({default: () => <span data-testid="system-prompt-section"/>}),
)
vi.mock(
    // eslint-disable-next-line max-len -- path too long to wrap
    "@/pages/workbench/components/WorkbenchTweaksPanel/components/TweaksConnectionsSection/TweaksConnectionsSection.tsx",
    () => ({default: () => <span data-testid="connections-section"/>}),
)
vi.mock(
    "@/pages/workbench/components/WorkbenchTweaksPanel/components/TweaksThemeSection/TweaksThemeSection.tsx",
    () => ({default: () => <span data-testid="theme-section"/>}),
)

function makeCtx(tweaksOpen: boolean) {
    return {tweaksOpen, tweaksSection: undefined, openTweaks: vi.fn(), closeTweaks: vi.fn()}
}

const base = {effectiveMode: "api" as const, status: "running", vaultId: "v1"}

describe("WorkbenchTweaksPanel", () => {
    it("keeps the header mounted but no section bodies while closed", () => {
        render(<WorkbenchTweaksPanel {...base} ctx={makeCtx(false)}/>)

        expect(screen.getByRole("heading", {name: "Tweaks"})).toBeInTheDocument()
        expect(screen.queryByTestId("theme-section")).not.toBeInTheDocument()
        expect(screen.queryByTestId("system-prompt-section")).not.toBeInTheDocument()
    })

    it("renders the section bodies when open, including System prompt in api mode", () => {
        render(<WorkbenchTweaksPanel {...base} effectiveMode="api" ctx={makeCtx(true)}/>)

        expect(screen.getByTestId("system-prompt-section")).toBeInTheDocument()
        expect(screen.getByTestId("theme-section")).toBeInTheDocument()
        expect(screen.getByTestId("connections-section")).toBeInTheDocument()
    })

    it("hides the System prompt section outside api mode", () => {
        render(<WorkbenchTweaksPanel {...base} effectiveMode="docker" ctx={makeCtx(true)}/>)

        expect(screen.queryByTestId("system-prompt-section")).not.toBeInTheDocument()
        expect(screen.getByTestId("theme-section")).toBeInTheDocument()
    })

    it("closes on the close button", () => {
        const ctx = makeCtx(true)
        render(<WorkbenchTweaksPanel {...base} ctx={ctx}/>)

        fireEvent.click(screen.getByRole("button", {name: "Close tweaks"}))

        expect(ctx.closeTweaks).toHaveBeenCalledTimes(1)
    })

    it("closes on Escape while open", () => {
        const ctx = makeCtx(true)
        render(<WorkbenchTweaksPanel {...base} ctx={ctx}/>)

        fireEvent.keyDown(document, {key: "Escape"})

        expect(ctx.closeTweaks).toHaveBeenCalledTimes(1)
    })

    it("ignores Escape while closed", () => {
        const ctx = makeCtx(false)
        render(<WorkbenchTweaksPanel {...base} ctx={ctx}/>)

        fireEvent.keyDown(document, {key: "Escape"})

        expect(ctx.closeTweaks).not.toHaveBeenCalled()
    })
})
