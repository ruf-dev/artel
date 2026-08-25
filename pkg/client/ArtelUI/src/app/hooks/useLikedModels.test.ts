import {act, renderHook} from "@testing-library/react"
import {beforeEach, describe, expect, it} from "vitest"

import {useLikedModels} from "@/app/hooks/useLikedModels.ts"

const STORAGE_KEY = "artel.likedOpenRouterModels"

beforeEach(() => {
    localStorage.clear()
})

describe("useLikedModels", () => {
    it("defaults to an empty liked list when localStorage is empty", () => {
        const {result} = renderHook(() => useLikedModels())

        expect(result.current.likedIds).toEqual([])
        expect(result.current.isLiked("openai/gpt-4")).toBe(false)
    })

    it("defaults to an empty liked list when localStorage holds invalid JSON", () => {
        localStorage.setItem(STORAGE_KEY, "{not valid json")

        const {result} = renderHook(() => useLikedModels())

        expect(result.current.likedIds).toEqual([])
    })

    it("defaults to an empty liked list when localStorage holds valid JSON that isn't an array", () => {
        localStorage.setItem(STORAGE_KEY, JSON.stringify({foo: "bar"}))

        const {result} = renderHook(() => useLikedModels())

        expect(result.current.likedIds).toEqual([])
    })

    it("filters out non-string entries from a malformed stored array", () => {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(["openai/gpt-4", 42, null, "anthropic/claude"]))

        const {result} = renderHook(() => useLikedModels())

        expect(result.current.likedIds).toEqual(["openai/gpt-4", "anthropic/claude"])
    })

    it("reads an existing liked list from localStorage", () => {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(["openai/gpt-4"]))

        const {result} = renderHook(() => useLikedModels())

        expect(result.current.likedIds).toEqual(["openai/gpt-4"])
        expect(result.current.isLiked("openai/gpt-4")).toBe(true)
        expect(result.current.isLiked("anthropic/claude")).toBe(false)
    })

    it("toggleLiked adds an id, persists it, and isLiked reflects it", () => {
        const {result} = renderHook(() => useLikedModels())

        act(() => {
            result.current.toggleLiked("openai/gpt-4")
        })

        expect(result.current.likedIds).toEqual(["openai/gpt-4"])
        expect(result.current.isLiked("openai/gpt-4")).toBe(true)
        expect(JSON.parse(localStorage.getItem(STORAGE_KEY)!)).toEqual(["openai/gpt-4"])
    })

    it("toggleLiked removes an already-liked id and persists the removal", () => {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(["openai/gpt-4", "anthropic/claude"]))
        const {result} = renderHook(() => useLikedModels())

        act(() => {
            result.current.toggleLiked("openai/gpt-4")
        })

        expect(result.current.likedIds).toEqual(["anthropic/claude"])
        expect(result.current.isLiked("openai/gpt-4")).toBe(false)
        expect(JSON.parse(localStorage.getItem(STORAGE_KEY)!)).toEqual(["anthropic/claude"])
    })

    it("persists across separate hook instances via localStorage", () => {
        const first = renderHook(() => useLikedModels())

        act(() => {
            first.result.current.toggleLiked("openai/gpt-4")
        })

        const second = renderHook(() => useLikedModels())

        expect(second.result.current.likedIds).toEqual(["openai/gpt-4"])
        expect(second.result.current.isLiked("openai/gpt-4")).toBe(true)
    })
})
