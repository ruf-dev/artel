import cls from "@/components/ProviderChip/ProviderChip.module.css"
import {cn} from "@/app/utils/cn"
import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import ProviderIcon from "@/components/ProviderIcon/ProviderIcon.tsx"

export default function ProviderChip({provider, variantClass, label}: {
    provider: ExternalProvider
    variantClass: string
    label: string
}) {
    return (
        <span
            className={cn(cls.ProviderChip, variantClass)}
            data-tooltip-id="root-tooltip"
            data-tooltip-content={`${label} connection`}
        >
            <span className={cls.ProviderChipIcon}><ProviderIcon provider={provider}/></span>
            {label}
        </span>
    )
}
