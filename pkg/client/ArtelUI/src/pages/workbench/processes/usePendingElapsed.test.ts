import {act, renderHook} from "@testing-library/react"
import {describe, expect, it, beforeEach, afterEach, vi} from "vitest"

import {usePendingElapsed, SLOW_AFTER_MS, STUCK_AFTER_MS} from "@/pages/workbench/processes/usePendingElapsed.ts"

describe("usePendingElapsed", () => {
    beforeEach(() => {
        vi.useFakeTimers()
    })

    afterEach(() => {
        vi.useRealTimers()
    })

    it("returns 'normal' when inactive", () => {
        const {result} = renderHook(() => usePendingElapsed(false, undefined))
        expect(result.current).toBe("normal")
    })

    it("returns 'normal' for the first 20s after becoming active", () => {
        const {result} = renderHook(() => usePendingElapsed(true, undefined))
        expect(result.current).toBe("normal")

        act(() => {
            vi.advanceTimersByTime(19000)
        })
        expect(result.current).toBe("normal")
    })

    it("returns 'slow' after 20s while active", () => {
        const {result} = renderHook(() => usePendingElapsed(true, undefined))
        expect(result.current).toBe("normal")

        act(() => {
            vi.advanceTimersByTime(SLOW_AFTER_MS)
        })
        expect(result.current).toBe("slow")
    })

    it("returns 'stuck' after 75s while active", () => {
        const {result} = renderHook(() => usePendingElapsed(true, undefined))
        expect(result.current).toBe("normal")

        act(() => {
            vi.advanceTimersByTime(STUCK_AFTER_MS)
        })
        expect(result.current).toBe("stuck")
    })

    it("resets to 'normal' when resetKey changes while active", () => {
        const {result, rerender} = renderHook(
            ({resetKey}: {resetKey: unknown}) => usePendingElapsed(true, resetKey),
            {initialProps: {resetKey: "key1"}},
        )

        act(() => {
            vi.advanceTimersByTime(SLOW_AFTER_MS)
        })
        expect(result.current).toBe("slow")

        act(() => {
            rerender({resetKey: "key2"})
        })
        expect(result.current).toBe("normal")
    })

    it("reverts to 'normal' when active becomes false", () => {
        const {result, rerender} = renderHook(
            ({active}: {active: boolean}) => usePendingElapsed(active, undefined),
            {initialProps: {active: true}},
        )

        act(() => {
            vi.advanceTimersByTime(SLOW_AFTER_MS)
        })
        expect(result.current).toBe("slow")

        act(() => {
            rerender({active: false})
        })
        expect(result.current).toBe("normal")
    })
})
