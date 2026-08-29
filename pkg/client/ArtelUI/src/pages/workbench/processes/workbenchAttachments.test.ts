import {describe, expect, it} from "vitest"

import {stripAttachmentsPreamble, withAttachmentsPreamble} from "@/pages/workbench/processes/workbenchAttachments.ts"

describe("withAttachmentsPreamble", () => {
    it("returns the text unchanged when no files are attached", () => {
        expect(withAttachmentsPreamble("do the thing", [])).toBe("do the thing")
    })

    it("prepends a single preamble line naming the attached files, then a blank line", () => {
        expect(withAttachmentsPreamble("summarise these", ["notes/a.md", "b.md"]))
            .toBe("[Attached vault files: notes/a.md, b.md]\n\nsummarise these")
    })

    it("still builds the preamble when the user text is empty", () => {
        expect(withAttachmentsPreamble("", ["a.md"])).toBe("[Attached vault files: a.md]\n\n")
    })
})

describe("stripAttachmentsPreamble", () => {
    it("returns the text unchanged when no files are attached", () => {
        expect(stripAttachmentsPreamble("[Attached vault files: a.md]\n\ndo the thing", []))
            .toBe("[Attached vault files: a.md]\n\ndo the thing")
    })

    it("round-trips with withAttachmentsPreamble for a single attached file", () => {
        const text = "summarise these"
        const paths = ["notes/a.md"]
        expect(stripAttachmentsPreamble(withAttachmentsPreamble(text, paths), paths)).toBe(text)
    })

    it("round-trips with withAttachmentsPreamble for multiple attached files", () => {
        const text = "compare these two"
        const paths = ["notes/a.md", "b.md"]
        expect(stripAttachmentsPreamble(withAttachmentsPreamble(text, paths), paths)).toBe(text)
    })

    it("returns the text unchanged when it doesn't start with the reconstructed preamble", () => {
        expect(stripAttachmentsPreamble("legacy message with no preamble", ["a.md"]))
            .toBe("legacy message with no preamble")
    })
})
