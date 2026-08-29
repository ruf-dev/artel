import {describe, expect, it, vi} from "vitest"
import {render, screen} from "@testing-library/react"

import TopbarLeft from "@/pages/workbench/components/WorkbenchTopbar/components/TopbarLeft/TopbarLeft.tsx"

interface MockModelSwitcherProps {
    models: string[]
    value: string
    isLoading?: boolean
    onChange: (model: string) => void
    disabled?: boolean
    placeholder?: string
    needsAttention?: boolean
}

let lastModelSwitcherProps: MockModelSwitcherProps | undefined

vi.mock("@/pages/workbench/components/ModelSwitcher/ModelSwitcher.tsx", () => ({
    default: (props: MockModelSwitcherProps) => {
        lastModelSwitcherProps = props
        return (
            <span
                data-testid="model-switcher"
                data-needs-attention={String(Boolean(props.needsAttention))}
            />
        )
    },
}))

vi.mock("@/pages/workbench/components/WorkbenchStatusBadge/WorkbenchStatusBadge.tsx", () => ({
    default: () => <span data-testid="status-badge"/>,
}))

describe("TopbarLeft", () => {
    const model = {models: ["a/b"], value: "a/b", isLoading: false, onChange: vi.fn()}

    it("flags needsAttention when in api mode with no model selected and not loading", () => {
        render(
            <TopbarLeft
                effectiveMode="api"
                exists
                status="running"
                model={{models: [], value: "", isLoading: false, onChange: vi.fn()}}
            />,
        )

        expect(lastModelSwitcherProps?.needsAttention).toBe(true)
    })

    it("does not flag needsAttention while models are still loading", () => {
        render(
            <TopbarLeft
                effectiveMode="api"
                exists
                status="running"
                model={{models: [], value: "", isLoading: true, onChange: vi.fn()}}
            />,
        )

        expect(lastModelSwitcherProps?.needsAttention).toBe(false)
    })

    it("does not flag needsAttention once a model is selected", () => {
        render(
            <TopbarLeft
                effectiveMode="api"
                exists
                status="running"
                model={{models: ["a/b"], value: "a/b", isLoading: false, onChange: vi.fn()}}
            />,
        )

        expect(lastModelSwitcherProps?.needsAttention).toBe(false)
    })

    it("does not render an api-mode ModelSwitcher in docker mode", () => {
        const {rerender} = render(
            <TopbarLeft effectiveMode="docker" exists status="running" model={model}/>,
        )

        // Docker mode: should have exactly 2 ModelSwitchers (status badge first, then disabled Claude Code switcher)
        // and a status badge
        expect(screen.getByTestId("status-badge")).toBeInTheDocument()

        // Rerender in api mode: should have exactly 1 ModelSwitcher (the live one)
        rerender(<TopbarLeft effectiveMode="api" exists status="running" model={model}/>)
        // Should not have status badge in api mode
        expect(screen.queryByTestId("status-badge")).not.toBeInTheDocument()
    })

    it("passes model props to the api-mode ModelSwitcher", () => {
        const onChange = vi.fn()
        const models = ["a/b", "c/d"]
        render(
            <TopbarLeft
                effectiveMode="api"
                exists
                status="running"
                model={{models, value: "a/b", isLoading: false, onChange}}
            />,
        )

        expect(lastModelSwitcherProps?.models).toEqual(models)
        expect(lastModelSwitcherProps?.value).toBe("a/b")
        expect(lastModelSwitcherProps?.onChange).toBe(onChange)
    })
})
