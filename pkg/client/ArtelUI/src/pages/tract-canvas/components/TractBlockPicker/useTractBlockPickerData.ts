import {useMemo} from "react"

import {TractTool} from "@/processes/Tracts.ts"
import {MomCandidate} from "@/app/api/artel/mcp_keys.pb.ts"
import {BranchIcon} from "@/pages/tract-canvas/icons/BranchIcon/BranchIcon.tsx"

const BUILTIN_MCP = "artel"

export interface LogicOption {
    type: "condition"
    name: string
    desc: string
    Icon: (p: { className?: string }) => React.JSX.Element
}

const LOGIC_OPTIONS: LogicOption[] = [
    {type: "condition", name: "Condition", desc: "Route to a true/false branch", Icon: BranchIcon},
]

interface Args {
    tools: TractTool[]
    momCandidates: MomCandidate[]
    query: string
    connFilter: string | null
}

function rank(mcp: string, availableMcps: Set<string | undefined>) {
    return availableMcps.has(mcp) ? 0 : 1
}

export function useTractBlockPickerData(args: Args) {
    const q = args.query.trim().toLowerCase()

    const connectedCandidates = useMemo(
        () => args.momCandidates.filter(c => (c.connections?.length ?? 0) > 0),
        [args.momCandidates],
    )
    const availableMcps = useMemo(
        () => new Set(connectedCandidates.map(c => c.name)),
        [connectedCandidates],
    )

    const grouped = useMemo(() => {
        const acc: Record<string, TractTool[]> = {}
        for (const t of args.tools) {
            if (t.mcp === BUILTIN_MCP) continue
            if (args.connFilter && t.mcp !== args.connFilter) continue
            if (
                q && !`${t.mcp}.${t.tool}`.toLowerCase().includes(q) && !t.description?.toLowerCase().includes(q)
            ) continue
            ;(acc[t.mcp] ??= []).push(t)
        }
        return acc
    }, [args.tools, q, args.connFilter])

    const orderedGroups = useMemo(() => {
        return Object.entries(grouped).sort(([mcpA], [mcpB]) => {
            const diff = rank(mcpA, availableMcps) - rank(mcpB, availableMcps)
            return diff !== 0 ? diff : mcpA.localeCompare(mcpB)
        })
    }, [grouped, availableMcps])

    const logicOptions = args.connFilter
        ? []
        : LOGIC_OPTIONS.filter(o => !q || o.name.toLowerCase().includes(q) || o.desc.toLowerCase().includes(q))

    return {connectedCandidates, grouped, orderedGroups, logicOptions}
}
