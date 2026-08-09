import {FeatureOverrideState}
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/admin/components/ArtelUsersTab/components/ArtelUserDetailDialog/components/UserSubscriptionDialog/components/FeatureOverrideRow/FeatureOverrideRow.tsx"

export const FEATURES: {key: string, label: string}[] = [
    {key: "emails", label: "Emails"},
    {key: "task_trackers", label: "Task Trackers"},
    {key: "notes", label: "Notes"},
    {key: "spreadsheets", label: "Spreadsheets"},
    {key: "connectors", label: "Connectors"},
    {key: "tract", label: "Tract"},
    {key: "postgres", label: "PostgreSQL"},
]

export function bytesToMb(bytes: string | undefined): number {
    if (!bytes) return 0
    return Math.round(Number(bytes) / (1024 * 1024))
}

export function mbToBytes(mb: string): string {
    const n = Number(mb) || 0
    return String(Math.round(n * 1024 * 1024))
}

export function initFeatureStates(
    overrides: Record<string, boolean> | undefined,
): Record<string, FeatureOverrideState> {
    const result: Record<string, FeatureOverrideState> = {}
    for (const feature of FEATURES) {
        if (overrides && feature.key in overrides) {
            result[feature.key] = overrides[feature.key] ? "on" : "off"
        } else {
            result[feature.key] = "inherit"
        }
    }
    return result
}
