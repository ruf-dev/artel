// Shared line-icon set for the tract canvas UI (list page, canvas nodes, block picker) — keeps
// icon choices consistent across all three instead of each widget inventing its own SVGs.

export interface IconProps {
    className?: string
}

export const base = {
    viewBox: "0 0 24 24",
    fill: "none",
    stroke: "currentColor",
    strokeWidth: 1.6,
    strokeLinecap: "round" as const,
    strokeLinejoin: "round" as const,
}
