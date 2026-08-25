import {describe, expect, it, vi} from "vitest"
import {render, screen} from "@testing-library/react"

import SimpleChatTopBar from "@/pages/workbench/components/SimpleChatTopBar/SimpleChatTopBar.tsx"

vi.mock("@/pages/workbench/components/WorkbenchTopBarShell/WorkbenchTopBarShell.tsx", () => ({
    default: (props: {vaultName: string; actions?: React.ReactNode}) => (
        <div className="top-bar-shell" data-testid="top-bar-shell">
            <span>{props.vaultName}</span>
            {props.actions}
        </div>
    ),
}))

// eslint-disable-next-line max-len -- mock module path too long to wrap under 120 chars
vi.mock("@/pages/workbench/components/SimpleChatTopBar/components/SimpleChatToolbarActions/SimpleChatToolbarActions.tsx", () => ({
    default: (props: {currentModel: string}) => (
        <div className="toolbar-actions" data-testid="toolbar-actions">{props.currentModel}</div>
    ),
}))

describe("SimpleChatTopBar", () => {
    it("passes vaultName through to the shell and renders the actions into its actions slot", () => {
        render(
            <SimpleChatTopBar
                vaultName="Test Vault"
                models={["m1"]}
                currentModel="m1"
                onChangeModel={vi.fn()}
                onNewChat={vi.fn()}
                onToggleHistory={vi.fn()}
            />,
        )

        expect(screen.getByText("Test Vault")).toBeInTheDocument()
        expect(screen.getByTestId("toolbar-actions")).toHaveTextContent("m1")
    })
})
