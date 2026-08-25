/* eslint-disable @typescript-eslint/no-explicit-any, max-lines */
import {describe, expect, it} from "vitest"

import {
    buildLatestModelFlags,
    filterModelTree,
    groupModelsByFamily,
    sortLikedToTop,
} from "@/processes/groupModelsByFamily.ts"

describe("groupModelsByFamily", () => {
    it("groups models by provider prefix and keeps singletons flat", () => {
        const models = [
            "anthropic/claude-3-opus",
            "anthropic/claude-3-sonnet",
            "openai/gpt-4",
            "google/gemini-2-5-pro",
        ]

        const result = groupModelsByFamily(models)

        expect(result).toHaveLength(3)
        // First item is the group (anthropic with 2 models)
        expect(result[0]).toEqual({
            group: "anthropic",
            options: [
                {id: "anthropic/claude-3-opus", name: "claude-3-opus"},
                {id: "anthropic/claude-3-sonnet", name: "claude-3-sonnet"},
            ],
        })
        // Then singles (google, openai each with 1 model, sorted alphabetically)
        expect(result[1]).toEqual({id: "google/gemini-2-5-pro", name: "google/gemini-2-5-pro"})
        expect(result[2]).toEqual({id: "openai/gpt-4", name: "openai/gpt-4"})
    })

    it("treats single-model families as flat leaves at the bottom", () => {
        const models = [
            "anthropic/claude-3-opus",
            "anthropic/claude-3-sonnet",
            "openai/gpt-4",
            "google/gemini-2-5-pro",
        ]

        const result = groupModelsByFamily(models)

        // Groups first (anthropic has 2 models, google has 1 so goes to singles)
        expect(result[0]).toEqual({
            group: "anthropic",
            options: [
                {id: "anthropic/claude-3-opus", name: "claude-3-opus"},
                {id: "anthropic/claude-3-sonnet", name: "claude-3-sonnet"},
            ],
        })
        // Then singles (sorted alphabetically: google, then openai)
        expect(result[1]).toEqual({id: "google/gemini-2-5-pro", name: "google/gemini-2-5-pro"})
        expect(result[2]).toEqual({id: "openai/gpt-4", name: "openai/gpt-4"})
    })

    it("handles models with no slash (single-model families with no provider separator)", () => {
        const models = ["standalone-model", "anthropic/claude-3-opus"]

        const result = groupModelsByFamily(models)

        expect(result).toHaveLength(2)
        // Both are singletons (1 model each), so both are flat leaves, sorted alphabetically
        expect(result[0]).toEqual({id: "anthropic/claude-3-opus", name: "anthropic/claude-3-opus"})
        expect(result[1]).toEqual({id: "standalone-model", name: "standalone-model"})
    })

    it("handles models with multiple slashes, keeping everything after the first slash as the name", () => {
        const models = ["meta-llama/llama-3.1-8b-instruct", "meta-llama/llama-2-7b"]

        const result = groupModelsByFamily(models)

        expect(result).toHaveLength(1)
        expect(result[0]).toEqual({
            group: "meta-llama",
            options: [
                {id: "meta-llama/llama-3.1-8b-instruct", name: "llama-3.1-8b-instruct"},
                {id: "meta-llama/llama-2-7b", name: "llama-2-7b"},
            ],
        })
    })

    it("preserves full model ID in leaf ID even when name is stripped", () => {
        const models = ["provider/model-1", "provider/model-2"]

        const result = groupModelsByFamily(models)

        expect(result).toHaveLength(1)
        const group = result[0] as any
        expect(group.options[0].id).toBe("provider/model-1")
        expect(group.options[0].name).toBe("model-1")
        expect(group.options[1].id).toBe("provider/model-2")
        expect(group.options[1].name).toBe("model-2")
    })

    it("sorts groups alphabetically and singles alphabetically", () => {
        const models = [
            "zebra/model",
            "apple/model-1",
            "apple/model-2",
            "banana/model",
            "cherry/model-1",
            "cherry/model-2",
        ]

        const result = groupModelsByFamily(models)

        // Groups first (apple, cherry sorted)
        expect((result[0] as any).group).toBe("apple")
        expect((result[1] as any).group).toBe("cherry")
        // Singles at bottom (banana, zebra sorted)
        expect((result[2] as any).id).toBe("banana/model")
        expect((result[3] as any).id).toBe("zebra/model")
    })
})

describe("groupModelsByFamily: ~ latest-alias collapsing", () => {
    it("merges a ~-prefixed family with its non-prefixed counterpart into one group", () => {
        const models = [
            "~anthropic/claude-3-opus",
            "anthropic/claude-3-opus",
            "anthropic/claude-3-sonnet",
        ]

        const result = groupModelsByFamily(models)

        expect(result).toHaveLength(1)
        const group = result[0] as any
        // Merged group key is normalized (no leading ~)
        expect(group.group).toBe("anthropic")
        expect(group.options).toHaveLength(3)
        // Ids are preserved exactly as given, including the leading ~
        const ids = group.options.map((o: any) => o.id)
        expect(ids).toContain("~anthropic/claude-3-opus")
        expect(ids).toContain("anthropic/claude-3-opus")
        expect(ids).toContain("anthropic/claude-3-sonnet")
        // Names are stripped using each model's own prefix length, so the ~
        // never leaks into the displayed name
        const tildeOpt = group.options.find((o: any) => o.id === "~anthropic/claude-3-opus")
        expect(tildeOpt.name).toBe("claude-3-opus")
        const plainOpt = group.options.find((o: any) => o.id === "anthropic/claude-3-opus")
        expect(plainOpt.name).toBe("claude-3-opus")
    })

    it("merges an all-~-prefixed family into a normalized group key too", () => {
        const models = ["~x/model-a", "~x/model-b"]

        const result = groupModelsByFamily(models)

        expect(result).toHaveLength(1)
        const group = result[0] as any
        expect(group.group).toBe("x")
        expect(group.options).toEqual([
            {id: "~x/model-a", name: "model-a"},
            {id: "~x/model-b", name: "model-b"},
        ])
    })

    it("keeps a lone ~-prefixed model as a singleton leaf with its ~ preserved in id and name", () => {
        const models = ["~solo/model", "other/model"]

        const result = groupModelsByFamily(models)

        expect(result).toHaveLength(2)
        expect(result.map((r: any) => r.id).sort()).toEqual(["other/model", "~solo/model"].sort())
        const soloLeaf = result.find((r: any) => r.id === "~solo/model") as any
        expect(soloLeaf.name).toBe("~solo/model")
    })
})

describe("buildLatestModelFlags", () => {
    it("flags only models whose family prefix carries a leading ~", () => {
        const models = [
            "~anthropic/claude-3-opus",
            "anthropic/claude-3-opus",
            "anthropic/claude-3-sonnet",
            "openai/gpt-4",
        ]

        const flags = buildLatestModelFlags(models)

        expect(flags.get("~anthropic/claude-3-opus")).toBe(true)
        expect(flags.get("anthropic/claude-3-opus")).toBe(false)
        expect(flags.get("anthropic/claude-3-sonnet")).toBe(false)
        expect(flags.get("openai/gpt-4")).toBe(false)
    })

    it("keys the map by the original unmodified model id, ~ included", () => {
        const flags = buildLatestModelFlags(["~x/model-a"])

        expect(flags.has("~x/model-a")).toBe(true)
        expect(flags.has("x/model-a")).toBe(false)
    })

    it("flags a ~-prefixed model with no slash", () => {
        const flags = buildLatestModelFlags(["~standalone-model", "standalone-model-2"])

        expect(flags.get("~standalone-model")).toBe(true)
        expect(flags.get("standalone-model-2")).toBe(false)
    })
})

describe("filterModelTree", () => {
    it("returns all items unchanged for empty query", () => {
        const items = [
            {group: "anthropic", options: [{id: "anthropic/claude", name: "claude"}]},
            {id: "openai/gpt-4", name: "openai/gpt-4"},
        ]

        expect(filterModelTree(items, "")).toEqual(items)
        expect(filterModelTree(items, "  ")).toEqual(items)
    })

    it("case-insensitively matches group names and keeps entire group", () => {
        const items = [
            {group: "anthropic", options: [{id: "anthropic/claude", name: "claude"}]},
            {group: "openai", options: [{id: "openai/gpt-4", name: "gpt-4"}]},
        ]

        const result = filterModelTree(items, "anthropic")

        expect(result).toHaveLength(1)
        expect((result[0] as any).group).toBe("anthropic")
    })

    it("case-insensitively matches leaf names within groups and keeps only matching leaves", () => {
        const items = [
            {
                group: "anthropic",
                options: [
                    {id: "anthropic/claude-opus", name: "claude-opus"},
                    {id: "anthropic/claude-sonnet", name: "claude-sonnet"},
                ],
            },
        ]

        const result = filterModelTree(items, "opus")

        expect(result).toHaveLength(1)
        const group = result[0] as any
        expect(group.group).toBe("anthropic")
        expect(group.options).toHaveLength(1)
        expect(group.options[0].name).toBe("claude-opus")
    })

    it("case-insensitively matches leaf IDs as well as names", () => {
        const items = [
            {id: "openai/gpt-4", name: "gpt-4"},
            {id: "openai/gpt-3.5", name: "gpt-3.5"},
        ]

        const result = filterModelTree(items, "gpt-4")

        expect(result).toHaveLength(1)
        expect((result[0] as any).id).toBe("openai/gpt-4")
    })

    it("drops groups that end up with zero matching leaves", () => {
        const items = [
            {
                group: "anthropic",
                options: [{id: "anthropic/claude", name: "claude"}],
            },
            {
                group: "openai",
                options: [{id: "openai/gpt-4", name: "gpt-4"}],
            },
        ]

        const result = filterModelTree(items, "gpt")

        expect(result).toHaveLength(1)
        expect((result[0] as any).group).toBe("openai")
    })

    it("preserves original top-level order (groups then singles)", () => {
        const items = [
            {group: "zebra", options: [{id: "zebra/model", name: "model"}]},
            {group: "apple", options: [{id: "apple/model", name: "model"}]},
            {id: "single/one", name: "single/one"},
        ]

        const result = filterModelTree(items, "")

        // Should keep original order (zebra, apple, single/one)
        expect((result[0] as any).group).toBe("zebra")
        expect((result[1] as any).group).toBe("apple")
        expect((result[2] as any).id).toBe("single/one")
    })

    it("handles mixed case search", () => {
        const items = [
            {group: "anthropic", options: [{id: "anthropic/Claude-Opus", name: "Claude-Opus"}]},
        ]

        const result = filterModelTree(items, "CLAUDE")

        expect(result).toHaveLength(1)
        expect((result[0] as any).options).toHaveLength(1)
    })

    it("matches within a merged ~-family group by group name, and by leaf id with a leading ~", () => {
        const merged = groupModelsByFamily([
            "~anthropic/claude-3-opus",
            "anthropic/claude-3-sonnet",
        ])

        // Matching the normalized group name keeps the whole merged group, ~ leaf included
        const byGroupName = filterModelTree(merged, "anthropic")
        expect(byGroupName).toHaveLength(1)
        expect((byGroupName[0] as any).options).toHaveLength(2)

        // Matching a substring that only appears in the ~-prefixed id still finds that leaf
        const byTildeId = filterModelTree(merged, "opus")
        expect(byTildeId).toHaveLength(1)
        expect((byTildeId[0] as any).options).toEqual([{id: "~anthropic/claude-3-opus", name: "claude-3-opus"}])
    })
})

describe("sortLikedToTop", () => {
    it("returns items unchanged when likedIds is empty", () => {
        const items = groupModelsByFamily(["anthropic/claude-3-opus", "openai/gpt-4"])

        const result = sortLikedToTop(items, [])

        expect(result).toEqual(items)
    })

    it("hoists a liked flat leaf into a Liked group at the top", () => {
        const items = groupModelsByFamily([
            "anthropic/claude-3-opus",
            "anthropic/claude-3-sonnet",
            "openai/gpt-4",
        ])

        const result = sortLikedToTop(items, ["openai/gpt-4"])

        expect(result).toHaveLength(2)
        expect((result[0] as any).group).toBe("Liked")
        expect((result[0] as any).options).toEqual([{id: "openai/gpt-4", name: "openai/gpt-4"}])
        expect((result[1] as any).group).toBe("anthropic")
    })

    it("extracts a liked option nested inside a multi-model family group", () => {
        const items = groupModelsByFamily(["anthropic/claude-3-opus", "anthropic/claude-3-sonnet", "openai/gpt-4"])

        const result = sortLikedToTop(items, ["anthropic/claude-3-sonnet"])

        expect(result).toHaveLength(3)
        // First item: Liked group with the extracted option
        expect((result[0] as any).group).toBe("Liked")
        expect((result[0] as any).options).toHaveLength(1)
        expect((result[0] as any).options[0].id).toBe("anthropic/claude-3-sonnet")
        // Second item: anthropic group with remaining options
        expect((result[1] as any).group).toBe("anthropic")
        expect((result[1] as any).options).toHaveLength(1)
        expect((result[1] as any).options[0].id).toBe("anthropic/claude-3-opus")
        // Third item: openai flat leaf (unchanged)
        expect((result[2] as any).id).toBe("openai/gpt-4")
    })

    it("drops a family group that becomes empty after extracting its only (liked) option", () => {
        const items = groupModelsByFamily(["anthropic/claude-3-opus", "anthropic/claude-3-sonnet", "openai/gpt-4"])

        const result = sortLikedToTop(items, ["openai/gpt-4"])

        // Should not have an empty group; openai should only appear in the Liked section
        expect(result).toHaveLength(2)
        expect((result[0] as any).group).toBe("Liked")
        expect((result[1] as any).group).toBe("anthropic")
        // Verify no singleton openai leaf appears after the anthropic group
        const hasEmptyOpenaiGroup = result.some((item: any) => item.group === "openai" && item.options?.length === 0)
        expect(hasEmptyOpenaiGroup).toBe(false)
    })

    it("keeps a group that still has non-liked options after extraction", () => {
        const items = groupModelsByFamily(["anthropic/claude-3-opus", "anthropic/claude-3-sonnet"])

        const result = sortLikedToTop(items, ["anthropic/claude-3-sonnet"])

        expect(result).toHaveLength(2)
        expect((result[0] as any).group).toBe("Liked")
        expect((result[1] as any).group).toBe("anthropic")
        expect((result[1] as any).options).toHaveLength(1)
        expect((result[1] as any).options[0].id).toBe("anthropic/claude-3-opus")
    })

    it("returns items unchanged when likedIds contains an id not in items", () => {
        const items = groupModelsByFamily(["anthropic/claude-3-opus", "openai/gpt-4"])

        const result = sortLikedToTop(items, ["nonexistent/model"])

        expect(result).toEqual(items)
    })

    it("collects liked options from multiple family groups into one Liked group", () => {
        const items = groupModelsByFamily([
            "anthropic/claude-3-opus",
            "anthropic/claude-3-sonnet",
            "openai/gpt-4",
            "openai/gpt-3.5-turbo",
            "mistral/mistral-large",
        ])

        const result = sortLikedToTop(items, [
            "anthropic/claude-3-sonnet",
            "openai/gpt-4",
            "mistral/mistral-large",
        ])

        expect(result).toHaveLength(3)
        // First: Liked group with all three liked options
        expect((result[0] as any).group).toBe("Liked")
        expect((result[0] as any).options).toHaveLength(3)
        const likedIds = (result[0] as any).options.map((opt: any) => opt.id)
        expect(likedIds).toContain("anthropic/claude-3-sonnet")
        expect(likedIds).toContain("openai/gpt-4")
        expect(likedIds).toContain("mistral/mistral-large")
        // Second: anthropic with remaining option
        expect((result[1] as any).group).toBe("anthropic")
        expect((result[1] as any).options).toHaveLength(1)
        expect((result[1] as any).options[0].id).toBe("anthropic/claude-3-opus")
        // Third: openai with remaining option
        expect((result[2] as any).group).toBe("openai")
        expect((result[2] as any).options).toHaveLength(1)
        expect((result[2] as any).options[0].id).toBe("openai/gpt-3.5-turbo")
    })
})
