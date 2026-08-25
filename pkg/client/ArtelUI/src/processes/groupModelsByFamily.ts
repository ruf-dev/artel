import {DropdownItem, DropdownOption, DropdownOptionGroup} from "@vervstack/chures"

/**
 * Split a raw model ID into its family prefix (everything before the first `/`,
 * or the whole ID if there's no `/`) and whether that prefix carries a leading
 * `~` — OpenRouter's marker for the "latest" rolling alias of a family (e.g.
 * `~anthropic/claude-3-opus` is the rolling alias of `anthropic/claude-3-opus`).
 * `familyKey` is the prefix with any leading `~` stripped, used to merge a
 * `~`-prefixed family with its non-prefixed counterpart into one group.
 */
function splitModelId(model: string): {rawPrefix: string; isLatest: boolean; familyKey: string} {
    const slashIndex = model.indexOf("/")
    const rawPrefix = slashIndex === -1 ? model : model.slice(0, slashIndex)
    const isLatest = rawPrefix.startsWith("~")
    const familyKey = isLatest ? rawPrefix.slice(1) : rawPrefix
    return {rawPrefix, isLatest, familyKey}
}

/**
 * Maps each raw model ID (unmodified — a leading `~` on its family prefix, if
 * any, is preserved) to whether it's the "latest" rolling alias of its family.
 * `DropdownOption` is a fixed `{id, name}` shape with no room for an extra
 * flag, so a caller that needs to render "latest" styling/badges looks it up
 * here by the `id` it gets back from `groupModelsByFamily`'s tree — that id is
 * never stripped of its `~` (only the *grouping key* is normalized, see
 * `groupModelsByFamily` below), so it round-trips as-is through selection.
 */
export function buildLatestModelFlags(models: string[]): Map<string, boolean> {
    const flags = new Map<string, boolean>()
    for (const model of models) {
        flags.set(model, splitModelId(model).isLatest)
    }
    return flags
}

/**
 * Group raw OpenRouter model IDs by their provider prefix (before the first `/`).
 * Families with 2+ models become a DropdownOptionGroup; singleton families become flat leaves.
 * Groups sort alphabetically first, then singleton leaves sort alphabetically at the bottom.
 *
 * A family prefix may carry a leading `~` (OpenRouter's "latest" rolling-alias
 * marker, e.g. `~anthropic/claude-3-opus`). The grouping key strips that `~` so
 * `~anthropic/...` and `anthropic/...` merge into a single `anthropic` group
 * instead of forming two sibling groups. Each option's `id` is kept exactly as
 * given (leading `~` and all) — only the group key/display name is normalized —
 * so ids still round-trip correctly as the `value` through existing selection
 * logic. Use `buildLatestModelFlags` to tell which ids came from a `~`-prefixed
 * family.
 */
export function groupModelsByFamily(models: string[]): DropdownItem[] {
    // Build a map of family key -> models. The key has any leading `~` on the
    // prefix stripped (see splitModelId), so a `~`-prefixed family merges with
    // its non-prefixed counterpart.
    const familyMap = new Map<string, string[]>()

    for (const model of models) {
        const {familyKey} = splitModelId(model)

        if (!familyMap.has(familyKey)) {
            familyMap.set(familyKey, [])
        }
        familyMap.get(familyKey)!.push(model)
    }

    const groups: DropdownOptionGroup[] = []
    const singles: DropdownOption[] = []

    // Separate into groups (2+ models) and singles (1 model)
    for (const [familyKey, modelsInFamily] of familyMap) {
        if (modelsInFamily.length >= 2) {
            groups.push({
                group: familyKey,
                options: modelsInFamily.map(m => {
                    // Strip using this model's own raw prefix length (which may
                    // include a leading `~`), not the normalized familyKey — so
                    // the displayed name never carries a stray `~` or gets
                    // sliced one character short/long.
                    const {rawPrefix} = splitModelId(m)
                    return {
                        id: m,
                        name: m.slice(rawPrefix.length + 1),
                    }
                }),
            })
        } else {
            // Singleton: the full model ID becomes the label
            singles.push({
                id: modelsInFamily[0],
                name: modelsInFamily[0],
            })
        }
    }

    // Sort groups by name, singles by name, then combine
    groups.sort((a, b) => a.group.localeCompare(b.group))
    singles.sort((a, b) => {
        const aId = typeof a === "string" ? a : a.id
        const bId = typeof b === "string" ? b : b.id
        return aId.localeCompare(bId)
    })

    return [...groups, ...singles]
}

/**
 * Extract a DropdownOption's id, handling both the plain-string and `{id, name}`
 * object variants. Mirrors chures' own internal `getOptionId` (Dropdown.types.ts),
 * which isn't part of the package's public exports, so callers that need this
 * outside chures (building a `renderOption` callback, looking up latest/liked
 * flags by id) share this small re-implementation instead of each hand-rolling
 * their own `typeof opt === "string" ? opt : opt.id` inline.
 */
export function getDropdownOptionId(opt: DropdownOption): string {
    return typeof opt === "string" ? opt : opt.id
}

/**
 * Extract a model's family string for icon lookup, e.g. "anthropic" from both
 * "anthropic/claude-3-opus" and "~anthropic/claude-3-opus" (leading `~` stripped,
 * same normalization as `groupModelsByFamily`'s grouping key). Returns the whole
 * id when there's no "/" separator.
 */
export function getModelFamily(id: string): string {
    return splitModelId(id).familyKey
}

function isGroupItem(item: DropdownItem): item is DropdownOptionGroup {
    return typeof item === "object" && item !== null && "group" in item && "options" in item
}

/**
 * Reorders a grouped `DropdownItem[]` tree so every model whose id is in
 * `likedIds` floats to a synthetic "Liked" group at the very top — including
 * models nested inside a family group (e.g. `anthropic/…`), not just
 * standalone top-level leaves. A liked model is extracted out of its
 * original family group; if that leaves the group empty, the group itself
 * is dropped (mirrors `filterModelTree`'s empty-group handling below).
 * Traversal order (groups first, then singles — same order
 * `groupModelsByFamily` already returns) determines the "Liked" group's
 * internal order. If none of `likedIds` match anything currently in
 * `items` (e.g. a stale localStorage entry for a model no longer in the
 * OpenRouter list), returns `items` unchanged rather than inserting an
 * empty "Liked" group header.
 *
 * Duplicate-id note: this always extracts-and-removes rather than
 * duplicating a liked option in both the "Liked" group and its original
 * family group. Reading chures' bundled Dropdown implementation confirmed
 * duplicate ids across two different groups would not actually break
 * anything (React keys and selection lookups are scoped per-group, not
 * globally) — extraction is still preferred purely because showing the same
 * model twice in one list is confusing UX.
 */
export function sortLikedToTop(items: DropdownItem[], likedIds: string[]): DropdownItem[] {
    if (likedIds.length === 0) {
        return items
    }

    const likedOptions: DropdownOption[] = []
    const rest: DropdownItem[] = []

    for (const item of items) {
        if (isGroupItem(item)) {
            const liked = item.options.filter(opt => likedIds.includes(getDropdownOptionId(opt)))
            const remaining = item.options.filter(opt => !likedIds.includes(getDropdownOptionId(opt)))
            likedOptions.push(...liked)
            if (remaining.length > 0) {
                rest.push({group: item.group, options: remaining})
            }
        } else if (likedIds.includes(getDropdownOptionId(item))) {
            likedOptions.push(item)
        } else {
            rest.push(item)
        }
    }

    if (likedOptions.length === 0) {
        return items
    }

    return [{group: "Liked", options: likedOptions}, ...rest]
}

/**
 * Filter a grouped model tree by search query (case-insensitive substring match).
 * Empty query returns items unchanged.
 * For groups, keeps the whole group if its name matches, otherwise keeps only matching leaves.
 * Drops groups that end up with zero leaves after filtering.
 * Preserves original top-level order (groups then singles, already sorted in groupModelsByFamily).
 */
export function filterModelTree(items: DropdownItem[], query: string): DropdownItem[] {
    const trimmedQuery = query.trim()

    // Empty query: return all items
    if (!trimmedQuery) {
        return items
    }

    const lowerQuery = trimmedQuery.toLowerCase()
    const result: Array<DropdownOption | DropdownOptionGroup> = []

    for (const item of items) {
        // Check if it's a group (has group and options properties)
        if (typeof item === "object" && item !== null && "group" in item && "options" in item) {
            const group = item as DropdownOptionGroup
            // Check if group name matches
            if (group.group.toLowerCase().includes(lowerQuery)) {
                // Keep entire group
                result.push(group)
            } else {
                // Filter group's options
                const filtered = group.options.filter(opt => {
                    const optName = typeof opt === "string" ? opt : opt.name
                    const optId = typeof opt === "string" ? opt : opt.id
                    return optName.toLowerCase().includes(lowerQuery) || optId.toLowerCase().includes(lowerQuery)
                })
                // Only keep group if it has matching leaves
                if (filtered.length > 0) {
                    result.push({
                        group: group.group,
                        options: filtered,
                    })
                }
            }
        } else {
            // It's a flat leaf
            const leaf = item as DropdownOption
            const leafName = typeof leaf === "string" ? leaf : leaf.name
            const leafId = typeof leaf === "string" ? leaf : leaf.id
            if (leafName.toLowerCase().includes(lowerQuery) || leafId.toLowerCase().includes(lowerQuery)) {
                result.push(leaf)
            }
        }
    }

    return result as DropdownItem[]
}
