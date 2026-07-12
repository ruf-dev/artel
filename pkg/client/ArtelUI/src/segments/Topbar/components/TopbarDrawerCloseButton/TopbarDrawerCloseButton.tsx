import {Button} from "@vervstack/chures"

import TopbarCloseIcon from "@/segments/Topbar/components/TopbarCloseIcon/TopbarCloseIcon.tsx"
import cls from "@/segments/Topbar/components/TopbarDrawerCloseButton/TopbarDrawerCloseButton.module.css"

interface TopbarDrawerCloseButtonProps {
    onClose: () => void
}

export default function TopbarDrawerCloseButton({onClose}: TopbarDrawerCloseButtonProps) {
    return (
        <div className={cls.TopbarDrawerCloseButtonContainer}>
            <Button
                variant="ghost"
                className={cls.CloseBtn}
                onClick={onClose}
                aria-label="Close navigation">
                <TopbarCloseIcon/>
            </Button>
        </div>
    )
}
