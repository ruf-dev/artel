import {beforeEach, describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import TweaksConnectionsSection
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/workbench/components/WorkbenchTweaksPanel/components/TweaksConnectionsSection/TweaksConnectionsSection.tsx"

const OpenDialog = vi.fn()
const fetch = vi.fn()
let connections: Array<Record<string, unknown>>

vi.mock("@/app/hooks/ExternalConnections.ts", () => ({
    useExternalConnections: () => ({connections, fetch}),
}))

vi.mock("@/app/hooks/Dialog", () => ({
    useDialog: () => ({OpenDialog}),
}))

vi.mock("@/dialogs/ManageOpenRouterDialog/ManageOpenRouterDialog.tsx", () => ({
    default: () => <span data-testid="manage-openrouter-dialog"/>,
}))

describe("TweaksConnectionsSection", () => {
    beforeEach(() => {
        OpenDialog.mockReset()
        fetch.mockReset().mockResolvedValue(undefined)
        connections = [{
            provider: ExternalProvider.EXTERNAL_PROVIDER_OPENROUTER,
            generic: {fields: {key_preview: "sk-or-****4f2a"}},
        }]
    })

    it("renders the OpenRouter masked key and opens the manage dialog", () => {
        render(<TweaksConnectionsSection effectiveMode="api" status="running"/>)

        expect(screen.getByText("sk-or-****4f2a")).toBeInTheDocument()

        fireEvent.click(screen.getByRole("button", {name: "Manage"}))

        expect(OpenDialog).toHaveBeenCalledTimes(1)
    })

    it("shows the container row with its status in docker mode", () => {
        render(<TweaksConnectionsSection effectiveMode="docker" status="stopped"/>)

        expect(screen.getByText("Claude Code container")).toBeInTheDocument()
        expect(screen.getByText("stopped")).toBeInTheDocument()
    })

    it("offers Connect when OpenRouter is not linked", () => {
        connections = []
        render(<TweaksConnectionsSection effectiveMode="api" status="running"/>)

        fireEvent.click(screen.getByRole("button", {name: "Connect"}))

        expect(OpenDialog).toHaveBeenCalledTimes(1)
    })
})
