import {describe, expect, it} from "vitest"

import {
    buildFolderTree,
    getAncestorFolderPaths,
    getDirectNotes,
    getItemName,
    sortItemsByName,
} from "@/components/FileTree/fileTree.ts"

describe("buildFolderTree", () => {
    it("builds a nested tree from flat folder paths", () => {
        const tree = buildFolderTree(["A", "A/B", "C"])

        expect(tree.map(n => n.name)).toEqual(["A", "C"])
        const a = tree.find(n => n.path === "A")
        expect(a?.children.map(c => c.path)).toEqual(["A/B"])
        expect(tree.find(n => n.path === "C")?.children).toEqual([])
    })

    it("synthesises intermediate folders that are not listed explicitly", () => {
        const tree = buildFolderTree(["x/y/z"])

        expect(tree.map(n => n.path)).toEqual(["x"])
        expect(tree[0].children[0].path).toBe("x/y")
        expect(tree[0].children[0].children[0].path).toBe("x/y/z")
    })

    it("sorts siblings by name", () => {
        expect(buildFolderTree(["b", "a", "c"]).map(n => n.name)).toEqual(["a", "b", "c"])
    })
})

describe("getItemName", () => {
    it("returns the basename without a trailing .md", () => {
        expect(getItemName({path: "Projects/spec.md"})).toBe("spec")
    })

    it("keeps non-md extensions", () => {
        expect(getItemName({path: "a/b/notes.txt"})).toBe("notes.txt")
    })

    it("falls back to Untitled when the path is missing", () => {
        expect(getItemName({})).toBe("Untitled")
    })
})

describe("sortItemsByName", () => {
    it("orders items by display name, not raw path", () => {
        const items = [{path: "z/apple.md"}, {path: "a/banana.md"}]

        expect(sortItemsByName(items).map(i => i.path)).toEqual(["z/apple.md", "a/banana.md"])
    })

    it("does not mutate its input", () => {
        const items = [{path: "b.md"}, {path: "a.md"}]

        sortItemsByName(items)

        expect(items.map(i => i.path)).toEqual(["b.md", "a.md"])
    })
})

describe("getDirectNotes", () => {
    const items = [
        {path: "Projects/spec.md"},
        {path: "Projects/sub/deep.md"},
        {path: "Projects/aaa.md"},
        {path: "Other/x.md"},
    ]

    it("returns only the folder's immediate children, sorted by name", () => {
        expect(getDirectNotes("Projects", items).map(i => i.path))
            .toEqual(["Projects/aaa.md", "Projects/spec.md"])
    })

    it("excludes items nested in subfolders", () => {
        expect(getDirectNotes("Projects", items).some(i => i.path === "Projects/sub/deep.md"))
            .toBe(false)
    })
})

describe("getAncestorFolderPaths", () => {
    it("lists every ancestor folder of a note path", () => {
        expect(getAncestorFolderPaths("a/b/c/note.md")).toEqual(["a", "a/b", "a/b/c"])
    })

    it("returns nothing for a root-level path", () => {
        expect(getAncestorFolderPaths("note.md")).toEqual([])
    })
})
