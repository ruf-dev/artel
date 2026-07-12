import {NavLink} from "react-router-dom"

import {cn} from "@/app/utils/cn.ts"
import {NavIconProps} from "@/segments/Topbar/components/icons/iconTypes.ts"
import cls from "@/segments/Topbar/components/TopbarMobileNavLink/TopbarMobileNavLink.module.css"

interface TopbarMobileNavLinkProps {
    to: string
    end?: boolean
    Icon: (props: NavIconProps) => React.JSX.Element
    label: string
    disabled?: boolean
}

export default function TopbarMobileNavLink({to, end, Icon, label, disabled}: TopbarMobileNavLinkProps) {
    if (disabled) {
        return (
            <span className={cls.NavLinkDisabled}>
                <Icon className={cls.NavIcon}/>
                {label}
            </span>
        )
    }

    return (
        <NavLink to={to} end={end} className={({isActive}) => cn(cls.NavLink, isActive && cls.NavLinkActive)}>
            <Icon className={cls.NavIcon}/>
            {label}
        </NavLink>
    )
}
