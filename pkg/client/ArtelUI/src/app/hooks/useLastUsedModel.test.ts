import {act, renderHook} from "@testing-library/react"
import {beforeEach, describe, expect, it} from "vitest"

import {useLastUsedModel} from "@/app/hooks/useLastUsedModel.ts"

const STORAGE_KEY = "artel.lastUsedSimpleChatModel"

beforeEach(() => {
    localStorage.clear()
})

describe("useLastUsedModel", () => {
    it("defaults to undefined when localStorage is empty", () => {
        const {result} = renderHook(() => useLastUsedModel())

        expect(result.current.lastUsedModel).toBeUndefined()
    })

    it("reads an existing value from localStorage", () => {
        localStorage.setItem(STORAGE_KEY, "openai/gpt-4")

        const {result} = renderHook(() => useLastUsedModel())

        expect(result.current.lastUsedModel).toBe("openai/gpt-4")
    })

    it("setLastUsedModel updates state and persists to localStorage", () => {
        const {result} = renderHook(() => useLastUsedModel())

        act(() => {
            result.current.setLastUsedModel("anthropic/claude")
        })

        expect(result.current.lastUsedModel).toBe("anthropic/claude")
        expect(localStorage.getItem(STORAGE_KEY)).toBe("anthropic/claude")
    })

    it("persists across separate hook instances via localStorage", () => {
        const first = renderHook(() => useLastUsedModel())

        act(() => {
            first.result.current.setLastUsedModel("anthropic/claude")
        })

        const second = renderHook(() => useLastUsedModel())

        expect(second.result.current.lastUsedModel).toBe("anthropic/claude")
    })
})
