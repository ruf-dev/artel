import {Button} from "@vervstack/chures"

import LogoutIcon from "@/segments/Topbar/components/icons/LogoutIcon.tsx"
import cls from "@/segments/Topbar/components/TopbarDrawerLogoutButton/TopbarDrawerLogoutButton.module.css"

interface TopbarDrawerLogoutButtonProps {
    onLogout: () => void
}

export default function TopbarDrawerLogoutButton({onLogout}: TopbarDrawerLogoutButtonProps) {
    return (
        <div className={cls.TopbarDrawerLogoutButtonContainer}>
            <Button variant="ghost" className={cls.LogoutBtn} onClick={onLogout}>
                <LogoutIcon className={cls.LogoutIcon}/>
                Log out
            </Button>
        </div>
    )
}
