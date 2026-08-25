import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import SimpleChatToolbarActions from
    "@/pages/workbench/components/SimpleChatTopBar/components/SimpleChatToolbarActions/SimpleChatToolbarActions.tsx"

vi.mock("@/pages/workbench/components/SimpleChat/components/ModelSwitcher/ModelSwitcher.tsx", () => ({
    default: (props: {value: string; onChange: (model: string) => void}) => (
        <div className="model-switcher" data-testid="model-switcher" onClick={() => props.onChange("new-model")}>
            {props.value}
        </div>
    ),
}))

describe("SimpleChatToolbarActions", () => {
    const defaultProps = {
        models: ["m1", "m2"],
        currentModel: "m1",
        onChangeModel: vi.fn(),
        onNewChat: vi.fn(),
        onToggleHistory: vi.fn(),
    }

    it("renders the model switcher with the current model", () => {
        render(<SimpleChatToolbarActions {...defaultProps}/>)

        expect(screen.getByTestId("model-switcher")).toHaveTextContent("m1")
    })

    it("calls onChangeModel when the model switcher changes", () => {
        const onChangeModel = vi.fn()
        render(<SimpleChatToolbarActions {...defaultProps} onChangeModel={onChangeModel}/>)

        fireEvent.click(screen.getByTestId("model-switcher"))

        expect(onChangeModel).toHaveBeenCalledWith("new-model")
    })

    it("calls onNewChat when the new-chat button is clicked", () => {
        const onNewChat = vi.fn()
        render(<SimpleChatToolbarActions {...defaultProps} onNewChat={onNewChat}/>)

        fireEvent.click(screen.getByLabelText("New chat"))

        expect(onNewChat).toHaveBeenCalledTimes(1)
    })

    it("calls onToggleHistory when the history button is clicked", () => {
        const onToggleHistory = vi.fn()
        render(<SimpleChatToolbarActions {...defaultProps} onToggleHistory={onToggleHistory}/>)

        fireEvent.click(screen.getByLabelText("View chat history"))

        expect(onToggleHistory).toHaveBeenCalledTimes(1)
    })
})
