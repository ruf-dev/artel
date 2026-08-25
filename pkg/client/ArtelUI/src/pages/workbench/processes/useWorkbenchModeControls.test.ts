import {act, renderHook} from "@testing-library/react"
import {describe, expect, it, vi} from "vitest"

import {useWorkbenchModeControls} from "@/pages/workbench/processes/useWorkbenchModeControls.ts"

// Only the toggleSimpleHistory/closeSimpleHistory additions are covered here —
// added so WorkbenchPage.tsx doesn't need its own local closures for these,
// mirroring useWorkbenchPanelControls.ts's toggleHistory/closeHistory for docker.
describe("useWorkbenchModeControls", () => {
    function renderControls() {
        return renderHook(() => useWorkbenchModeControls({exists: false, handleCreateDocker: vi.fn()}))
    }

    it("toggleSimpleHistory flips simpleHistoryOpen", () => {
        const {result} = renderControls()

        expect(result.current.simpleHistoryOpen).toBe(false)

        act(() => result.current.toggleSimpleHistory())
        expect(result.current.simpleHistoryOpen).toBe(true)

        act(() => result.current.toggleSimpleHistory())
        expect(result.current.simpleHistoryOpen).toBe(false)
    })

    it("closeSimpleHistory closes it when a simple chat is selected", () => {
        const {result} = renderControls()

        act(() => result.current.handleSimpleChatCreated("chat-1"))
        act(() => result.current.toggleSimpleHistory())
        expect(result.current.simpleHistoryOpen).toBe(true)

        act(() => result.current.closeSimpleHistory())
        expect(result.current.simpleHistoryOpen).toBe(false)
    })

    it("closeSimpleHistory is a no-op when there is no active simple chat", () => {
        const {result} = renderControls()

        act(() => result.current.setSimpleHistoryOpen(true))
        act(() => result.current.closeSimpleHistory())

        expect(result.current.simpleHistoryOpen).toBe(true)
    })
})
