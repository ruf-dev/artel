import {describe, expect, it} from "vitest"

import {filterVaultTree} from
    "@/pages/workbench/components/WorkbenchSidebar/components/VaultPane/vaultPaneFilter.ts"

const folders = ["Projects", "Projects/Alpha", "Archive"]
const notes = [
    {path: "Projects/Alpha/spec.md"},
    {path: "Projects/roadmap.md"},
    {path: "Archive/old.md"},
    {path: "root.md"},
]

describe("filterVaultTree", () => {
    it("returns folders and notes untouched when the query is blank", () => {
        expect(filterVaultTree("   ", folders, notes)).toEqual({folders, notes})
    })

    it("keeps only notes whose path matches the query, case-insensitively", () => {
        const out = filterVaultTree("SPEC", folders, notes)

        expect(out.notes).toEqual([{path: "Projects/Alpha/spec.md"}])
    })

    it("rebuilds the folder list as the union of matched notes' ancestor folders", () => {
        const out = filterVaultTree("projects", folders, notes)

        expect(out.notes.map(n => n.path)).toEqual([
            "Projects/Alpha/spec.md",
            "Projects/roadmap.md",
        ])
        expect([...out.folders].sort()).toEqual(["Projects", "Projects/Alpha"])
    })

    it("yields empty folders/notes when nothing matches", () => {
        expect(filterVaultTree("zzz", folders, notes)).toEqual({folders: [], notes: []})
    })

    it("matches a root note and produces no folders for it", () => {
        const out = filterVaultTree("root", folders, notes)

        expect(out.notes).toEqual([{path: "root.md"}])
        expect(out.folders).toEqual([])
    })
})
