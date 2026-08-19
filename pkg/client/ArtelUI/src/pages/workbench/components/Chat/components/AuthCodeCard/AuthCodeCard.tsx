import cls from "@/pages/workbench/components/Chat/components/AuthCodeCard/AuthCodeCard.module.css"
import {cn} from "@/app/utils/cn.ts"

interface Props {
    resolved: boolean
}

export default function AuthCodeCard({resolved}: Props) {
    return (
        <div className={cn(cls.AuthCodeCardContainer, resolved && cls.Resolved)}>
            <p className={cls.Text}>
                {resolved
                    ? "Authentication code submitted."
                    : "Enter the authentication code below to continue."}
            </p>
        </div>
    )
}
