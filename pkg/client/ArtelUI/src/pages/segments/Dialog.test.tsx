import {afterEach, describe, expect, it} from "vitest"
import {act, cleanup, render, screen} from "@testing-library/react"

import Dialog from "@/pages/segments/Dialog.tsx"
import {useDialog} from "@/app/hooks/Dialog"

afterEach(() => {
    act(() => {
        useDialog.setState({children: null, closable: false, IsClickOffClosesDialog: true})
    })
    cleanup()
    document.querySelectorAll("[data-earlier-overlay]").forEach(n => n.remove())
})

describe("Dialog", () => {
    it("renders nothing while there are no children", () => {
        render(<Dialog/>)

        expect(document.querySelector("[data-dialog-root]")?.childElementCount ?? 0).toBe(0)
    })

    it("portals the dialog content once opened", () => {
        render(<Dialog/>)

        act(() => {
            useDialog.getState().OpenDialog(<span>hello dialog</span>)
        })

        expect(screen.getByText("hello dialog")).toBeInTheDocument()
        expect(document.querySelector("[data-dialog-root]")).not.toBeNull()
    })

    it("keeps its portal target as the last child of body, after overlays mounted earlier", () => {
        // Simulate an overlay (e.g. WorkbenchTweaksPanel) portaled to body first.
        const earlier = document.createElement("div")
        earlier.setAttribute("data-earlier-overlay", "")
        document.body.appendChild(earlier)

        render(<Dialog/>)

        act(() => {
            useDialog.getState().OpenDialog(<span>on top</span>)
        })

        expect(document.body.lastElementChild).toHaveAttribute("data-dialog-root")
    })
})
