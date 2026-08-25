import {describe, expect, it} from "vitest"

import {getModelFamily, sortLikedToTop} from "@/processes/groupModelsByFamily.ts"

describe("getModelFamily", () => {
    it("returns the same family for a ~-prefixed id and a non-prefixed id of the same family", () => {
        expect(getModelFamily("~anthropic/claude-3-opus")).toBe("anthropic")
        expect(getModelFamily("anthropic/claude-3-opus")).toBe("anthropic")
        expect(getModelFamily("~anthropic/claude-3-opus")).toBe(getModelFamily("anthropic/claude-3-opus"))
    })

    it("returns the whole id when there is no slash", () => {
        expect(getModelFamily("standalone-model")).toBe("standalone-model")
    })
})

describe("sortLikedToTop", () => {
    it("moves a liked flat leaf to the front", () => {
        const items = [
            {id: "openai/gpt-4", name: "openai/gpt-4"},
            {id: "google/gemini-2-5-pro", name: "google/gemini-2-5-pro"},
        ]

        const result = sortLikedToTop(items, ["google/gemini-2-5-pro"])

        expect(result).toEqual([
            {id: "google/gemini-2-5-pro", name: "google/gemini-2-5-pro"},
            {id: "openai/gpt-4", name: "openai/gpt-4"},
        ])
    })

    it("does not hoist a liked leaf out of its group", () => {
        const items = [
            {
                group: "anthropic",
                options: [
                    {id: "anthropic/claude-3-opus", name: "claude-3-opus"},
                    {id: "anthropic/claude-3-sonnet", name: "claude-3-sonnet"},
                ],
            },
            {id: "openai/gpt-4", name: "openai/gpt-4"},
        ]

        const result = sortLikedToTop(items, ["anthropic/claude-3-opus"])

        // Group stays in its original position, its options untouched.
        expect(result).toEqual(items)
    })

    it("is a no-op for empty likedIds", () => {
        const items = [
            {id: "openai/gpt-4", name: "openai/gpt-4"},
            {group: "anthropic", options: [{id: "anthropic/claude-3-opus", name: "claude-3-opus"}]},
        ]

        expect(sortLikedToTop(items, [])).toBe(items)
    })
})
