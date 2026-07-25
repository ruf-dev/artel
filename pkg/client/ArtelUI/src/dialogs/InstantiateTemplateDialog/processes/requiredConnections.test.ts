import {describe, expect, it} from "bun:test"

import type {TractStep} from "@/processes/Tracts.ts"
import {requiredConnections} from "@/dialogs/InstantiateTemplateDialog/processes/requiredConnections"

function assertResult(
    result: ReturnType<typeof requiredConnections>,
    expected: ReturnType<typeof requiredConnections>,
): void {
    expect(JSON.stringify(result)).toBe(JSON.stringify(expected))
}

describe("requiredConnections", () => {
    it("flat list with distinct mcp values → dedup, exclude artel, keep order", () => {
        const steps: TractStep[] = [
            {type: "action", id: "1", mcp: "gitlab"},
            {type: "action", id: "2", mcp: "trello"},
            {type: "action", id: "3", mcp: "gitlab"},
            {type: "action", id: "4", mcp: "artel"},
        ]
        assertResult(requiredConnections(steps), [
            {key: "gitlab", kind: "mom"},
            {key: "trello", kind: "mom"},
        ])
    })

    it("nested under condition.then/else → recursion discovered", () => {
        const steps: TractStep[] = [
            {
                type: "condition",
                id: "1",
                then: [{type: "action", id: "2", mcp: "gitlab"}],
                else: [{type: "action", id: "3", mcp: "trello"}],
            },
        ]
        assertResult(requiredConnections(steps), [
            {key: "gitlab", kind: "mom"},
            {key: "trello", kind: "mom"},
        ])
    })

    it("nested under parallel/group.steps → recursion discovered", () => {
        const steps: TractStep[] = [
            {
                type: "parallel",
                id: "1",
                steps: [
                    {type: "action", id: "2", mcp: "gitlab"},
                    {type: "action", id: "3", mcp: "trello"},
                ],
            },
        ]
        assertResult(requiredConnections(steps), [
            {key: "gitlab", kind: "mom"},
            {key: "trello", kind: "mom"},
        ])
    })

    it("single llm_call → includes anthropic requirement", () => {
        const steps: TractStep[] = [
            {type: "action", id: "1", mcp: "gitlab"},
            {type: "llm_call", id: "2", prompt: "test"},
        ]
        assertResult(requiredConnections(steps), [
            {key: "gitlab", kind: "mom"},
            {key: "anthropic", kind: "llm"},
        ])
    })

    it("multiple llm_call steps → dedup to one anthropic requirement", () => {
        const steps: TractStep[] = [
            {type: "llm_call", id: "1", prompt: "test1"},
            {type: "llm_call", id: "2", prompt: "test2"},
            {type: "llm_call", id: "3", prompt: "test3"},
        ]
        assertResult(requiredConnections(steps), [
            {key: "anthropic", kind: "llm"},
        ])
    })

    it("both mom actions and llm_call → both kinds present", () => {
        const steps: TractStep[] = [
            {type: "action", id: "1", mcp: "gitlab"},
            {
                type: "condition",
                id: "2",
                then: [{type: "llm_call", id: "3", prompt: "test"}],
            },
        ]
        assertResult(requiredConnections(steps), [
            {key: "gitlab", kind: "mom"},
            {key: "anthropic", kind: "llm"},
        ])
    })

    it("empty steps → no requirements", () => {
        assertResult(requiredConnections([]), [])
    })

    it("only builtin actions → no requirements", () => {
        const steps: TractStep[] = [
            {type: "action", id: "1", mcp: "artel"},
            {type: "action", id: "2", mcp: "artel"},
        ]
        assertResult(requiredConnections(steps), [])
    })
})
