import {Button} from "@vervstack/chures"

import {cn} from "@/app/utils/cn.ts"
import TopbarHamburgerIcon from "@/segments/Topbar/components/TopbarHamburgerIcon/TopbarHamburgerIcon.tsx"
import cls from "@/segments/Topbar/components/TopbarMobileTrigger/TopbarMobileTrigger.module.css"

interface TopbarMobileTriggerProps {
    open: boolean
    onClick: () => void
}

export default function TopbarMobileTrigger({open, onClick}: TopbarMobileTriggerProps) {
    return (
        <div className={cls.TopbarMobileTriggerContainer}>
            <Button
                variant="ghost"
                className={cn(cls.HamburgerBtn, open && cls.HamburgerBtnActive)}
                onClick={onClick}
                aria-label="Open navigation"
            >
                <TopbarHamburgerIcon/>
            </Button>
        </div>
    )
}
