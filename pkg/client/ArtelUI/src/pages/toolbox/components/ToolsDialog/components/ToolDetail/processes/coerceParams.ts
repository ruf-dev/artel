import {ToolParamDef} from "@/app/api/artel/mcp_keys.pb.ts"

export function coerceParams(defs: Record<string, ToolParamDef>, values: Record<string, string>): Record<string, unknown> {
    const out: Record<string, unknown> = {}
    for (const [name, def] of Object.entries(defs)) {
        const raw = values[name]
        if (raw === undefined || raw === "") continue
        if (def.integerParam) {
            const n = Number(raw)
            out[name] = Number.isNaN(n) ? raw : n
        } else {
            out[name] = raw
        }
    }
    return out
}
