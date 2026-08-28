import {describe, expect, it} from "vitest"

import {withAttachmentsPreamble} from "@/pages/workbench/processes/workbenchAttachments.ts"

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
