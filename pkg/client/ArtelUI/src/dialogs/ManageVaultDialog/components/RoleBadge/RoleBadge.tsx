import cls from "@/dialogs/ManageVaultDialog/components/RoleBadge/RoleBadge.module.css"
import {cn} from "@/app/utils/cn.ts"

interface Props {
    role: string
}

export default function RoleBadge({role}: Props) {
    return <span className={cn(cls.RoleBadgeContainer, cls[`Role_${role}`])}>{role}</span>
}
