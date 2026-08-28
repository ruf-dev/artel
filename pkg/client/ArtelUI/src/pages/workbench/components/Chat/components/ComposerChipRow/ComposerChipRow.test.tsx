import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import ComposerChipRow from "@/pages/workbench/components/Chat/components/ComposerChipRow/ComposerChipRow.tsx"

describe("ComposerChipRow", () => {
    it("opens the tweaks panel at the matching section per chip", () => {
        const onOpenTweaks = vi.fn()
        render(<ComposerChipRow onOpenTweaks={onOpenTweaks}/>)

        fireEvent.click(screen.getByRole("button", {name: "Connections"}))
        expect(onOpenTweaks).toHaveBeenLastCalledWith("connections")

        fireEvent.click(screen.getByRole("button", {name: "System prompt"}))
        expect(onOpenTweaks).toHaveBeenLastCalledWith("system")

        fireEvent.click(screen.getByRole("button", {name: "Max tokens"}))
        expect(onOpenTweaks).toHaveBeenLastCalledWith("tokens")

        fireEvent.click(screen.getByRole("button", {name: "Context"}))
        expect(onOpenTweaks).toHaveBeenLastCalledWith("context")
    })
})
