import {act, renderHook} from "@testing-library/react"
import {describe, expect, it} from "vitest"

import {useWorkbenchContext} from "@/pages/workbench/processes/workbenchContext.ts"

describe("useWorkbenchContext", () => {
    it("starts closed with no section", () => {
        const {result} = renderHook(() => useWorkbenchContext())

        expect(result.current.tweaksOpen).toBe(false)
        expect(result.current.tweaksSection).toBeUndefined()
    })

    it("openTweaks(section) opens the panel and records the section", () => {
        const {result} = renderHook(() => useWorkbenchContext())

        act(() => result.current.openTweaks("connections"))

        expect(result.current.tweaksOpen).toBe(true)
        expect(result.current.tweaksSection).toBe("connections")
    })

    it("openTweaks() with no section opens the panel and clears the section", () => {
        const {result} = renderHook(() => useWorkbenchContext())

        act(() => result.current.openTweaks("system"))
        act(() => result.current.openTweaks())

        expect(result.current.tweaksOpen).toBe(true)
        expect(result.current.tweaksSection).toBeUndefined()
    })

    it("closeTweaks() closes the panel", () => {
        const {result} = renderHook(() => useWorkbenchContext())

        act(() => result.current.openTweaks("theme"))
        act(() => result.current.closeTweaks())

        expect(result.current.tweaksOpen).toBe(false)
    })

    it("toggles open then closed via openTweaks/closeTweaks", () => {
        const {result} = renderHook(() => useWorkbenchContext())

        expect(result.current.tweaksOpen).toBe(false)
        act(() => result.current.openTweaks())
        expect(result.current.tweaksOpen).toBe(true)
        act(() => result.current.closeTweaks())
        expect(result.current.tweaksOpen).toBe(false)
    })
})
