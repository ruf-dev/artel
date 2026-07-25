// Non-component helpers split out from the icons/ folder — react-refresh/only-export-components
// requires a component file to only export components.

import {ChatIcon} from "@/pages/tract-canvas/icons/ChatIcon/ChatIcon.tsx"
import {EditIcon} from "@/pages/tract-canvas/icons/EditIcon/EditIcon.tsx"
import {EnvelopeIcon} from "@/pages/tract-canvas/icons/EnvelopeIcon/EnvelopeIcon.tsx"
import {GlobeIcon} from "@/pages/tract-canvas/icons/GlobeIcon/GlobeIcon.tsx"
import {IconProps} from "@/pages/tract-canvas/icons/iconTypes.ts"
import {PaperPlaneIcon} from "@/pages/tract-canvas/icons/PaperPlaneIcon/PaperPlaneIcon.tsx"
import {PlusIcon} from "@/pages/tract-canvas/icons/PlusIcon/PlusIcon.tsx"
import {SearchIcon} from "@/pages/tract-canvas/icons/SearchIcon/SearchIcon.tsx"
import {TrashIcon} from "@/pages/tract-canvas/icons/TrashIcon/TrashIcon.tsx"
import {WrenchIcon} from "@/pages/tract-canvas/icons/WrenchIcon/WrenchIcon.tsx"

/** iconForTool picks a reasonable icon for an action step from its mcp/tool name — a heuristic
 * since real MoM tools don't carry an icon/category of their own. */
export function iconForTool(mcp: string, tool: string): (props: IconProps) => React.JSX.Element {
    const t = tool.toLowerCase()
    if (mcp === "telegram" || t.includes("telegram")) return PaperPlaneIcon
    if (t.includes("mail") || t.includes("email")) return EnvelopeIcon
    if (t.includes("chat") || t.includes("message") || t.includes("reply")) return ChatIcon
    if (t.includes("get") || t.includes("list") || t.includes("search") || t.includes("find")) return SearchIcon
    if (t.includes("update") || t.includes("edit") || t.includes("set")) return EditIcon
    if (t.includes("create") || t.includes("add") || t.includes("insert")) return PlusIcon
    if (t.includes("delete") || t.includes("remove")) return TrashIcon
    if (t.includes("http") || t.includes("request") || t.includes("get_url") || t.includes("post")) return GlobeIcon
    return WrenchIcon
}

export interface StepColor {
    bg: string
    border: string
    fg: string
}

/** colorForKind mirrors the mockup's color law: coral = builtin/vault ops, blue = external
 * integrations, amber = logic/condition, purple = structural (parallel/group), green = script. */
export function colorForKind(
    kind: "trigger" | "action" | "condition" | "parallel" | "group" | "script" | "llm_call",
    mcp?: string,
): StepColor {
    if (kind === "condition") {
        return {bg: "var(--amber-dim)", border: "rgba(245,158,11,0.28)", fg: "var(--color-warning)"}
    }
    if (kind === "parallel" || kind === "group") {
        return {bg: "rgba(167,139,250,0.12)", border: "rgba(167,139,250,0.28)", fg: "#a78bfa"}
    }
    if (kind === "script") {
        return {bg: "rgba(52,211,153,0.12)", border: "rgba(52,211,153,0.28)", fg: "#34d399"}
    }
    if (kind === "llm_call") {
        return {bg: "rgba(45,212,191,0.12)", border: "rgba(45,212,191,0.28)", fg: "#2dd4bf"}
    }
    if (kind === "trigger" || mcp === "artel" || !mcp) {
        return {bg: "var(--coral-dim)", border: "var(--coral-border)", fg: "var(--coral)"}
    }
    return {bg: "rgba(96,165,250,0.12)", border: "rgba(96,165,250,0.28)", fg: "#60a5fa"}
}
