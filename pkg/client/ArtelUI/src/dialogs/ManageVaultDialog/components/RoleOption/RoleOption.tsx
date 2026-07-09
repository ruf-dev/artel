import {Button} from "@vervstack/chures"

import cls from "@/dialogs/ManageVaultDialog/components/RoleOption/RoleOption.module.css"
import {cn} from "@/app/utils/cn.ts"

interface Props {
    selected: boolean
    label: string
    desc: string
    onSelect: () => void
}

export default function RoleOption({selected, label, desc, onSelect}: Props) {
    return (
        <Button
            className={cn(cls.RoleOptionContainer, selected && cls.RoleOptionSelected)}
            onClick={onSelect}
        >
            <span className={cls.RoleOptionLabel}>{label}</span>
            <span className={cls.RoleOptionDesc}>{desc}</span>
        </Button>
    )
}
