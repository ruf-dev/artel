import {act, renderHook} from "@testing-library/react"
import {describe, expect, it} from "vitest"

import {useWorkbenchAttachments} from "@/pages/workbench/processes/useWorkbenchAttachments.ts"

describe("useWorkbenchAttachments", () => {
    it("starts with no attached paths", () => {
        const {result} = renderHook(() => useWorkbenchAttachments())

        expect(result.current.paths).toEqual([])
    })

    it("toggle adds a path when absent and removes it when present", () => {
        const {result} = renderHook(() => useWorkbenchAttachments())

        act(() => result.current.toggle("a.md"))
        expect(result.current.paths).toEqual(["a.md"])

        act(() => result.current.toggle("b.md"))
        expect(result.current.paths).toEqual(["a.md", "b.md"])

        act(() => result.current.toggle("a.md"))
        expect(result.current.paths).toEqual(["b.md"])
    })

    it("remove drops the given path and is a no-op for an unknown one", () => {
        const {result} = renderHook(() => useWorkbenchAttachments())

        act(() => result.current.toggle("a.md"))
        act(() => result.current.toggle("b.md"))

        act(() => result.current.remove("a.md"))
        expect(result.current.paths).toEqual(["b.md"])

        act(() => result.current.remove("missing.md"))
        expect(result.current.paths).toEqual(["b.md"])
    })

    it("clear empties the set", () => {
        const {result} = renderHook(() => useWorkbenchAttachments())

        act(() => result.current.toggle("a.md"))
        act(() => result.current.toggle("b.md"))
        act(() => result.current.clear())

        expect(result.current.paths).toEqual([])
    })

    it("keeps stable callback identities across renders", () => {
        const {result, rerender} = renderHook(() => useWorkbenchAttachments())

        const first = {
            toggle: result.current.toggle,
            remove: result.current.remove,
            clear: result.current.clear,
        }

        act(() => result.current.toggle("a.md"))
        rerender()

        expect(result.current.toggle).toBe(first.toggle)
        expect(result.current.remove).toBe(first.remove)
        expect(result.current.clear).toBe(first.clear)
    })
})
