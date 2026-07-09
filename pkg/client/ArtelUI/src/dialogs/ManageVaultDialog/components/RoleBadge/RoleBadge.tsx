import cls from "@/dialogs/ManageVaultDialog/components/RoleBadge/RoleBadge.module.css"
import {cn} from "@/app/utils/cn.ts"

interface Props {
    role: string
}

export default function RoleBadge({role}: Props) {
    const roleClass = `Role${role.charAt(0).toUpperCase()}${role.slice(1)}`
    return <span className={cn(cls.RoleBadgeContainer, cls[roleClass])}>{role}</span>
}
