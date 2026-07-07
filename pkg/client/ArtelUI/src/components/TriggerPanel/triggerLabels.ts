import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import {Trigger, TriggerSource} from "@/processes/Tracts.ts"

// -- Category/provider display helpers (shared by TriggerPanel's TriggerRow and AddTriggerDialog's source picker) --

const PROVIDER_ENUM_BY_KEY: Record<string, ExternalProvider> = {
    gitlab: ExternalProvider.EXTERNAL_PROVIDER_GITLAB,
}

export function categoryLabel(category: string): string {
    switch (category) {
        case "gitlab":
            return "GitLab hooks"
        case "generic":
            return "Generic webhook"
        default:
            return category
    }
}

export function providerLabel(provider: string): string {
    if (provider === "gitlab") return "GitLab"
    return provider ? provider.charAt(0).toUpperCase() + provider.slice(1) : ""
}

export function providerEnumFor(provider: string): ExternalProvider | undefined {
    return PROVIDER_ENUM_BY_KEY[provider]
}

// One chip that says what a trigger is about — a webhook shows its preset label (e.g.
// "GitLab push"), a manual trigger just says "Manual" since its source is a schema-only
// placeholder ("generic"), not a meaningful preset.
export function triggerChipLabel(trigger: Trigger, preset: TriggerSource | undefined): string {
    if (trigger.kind === "manual") return "Manual"
    return preset?.label || trigger.source || trigger.kind
}
