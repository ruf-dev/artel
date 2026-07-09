import cls from "@/components/GenericChip/GenericChip.module.css"
import ProviderIcon from "@/components/ProviderIcon/ProviderIcon.tsx"

export default function GenericChip({mcpName}: {mcpName?: string}) {
    return (
        <span
            className={cls.GenericChip}
            data-tooltip-id="root-tooltip"
            data-tooltip-content={`${mcpName} connection`}
        >
            {/* TODO: no dedicated icon for this connection type yet - falls back to a generic glyph */}
            <span className={cls.ProviderChipIcon}><ProviderIcon/></span>
            {mcpName}
        </span>
    )
}
