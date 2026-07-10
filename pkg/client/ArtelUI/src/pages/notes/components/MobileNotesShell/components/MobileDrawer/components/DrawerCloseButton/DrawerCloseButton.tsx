import {Button} from "@vervstack/chures"

import ChevronLeftIcon from "@/pages/notes/components/icons/ChevronLeftIcon.tsx"
import cls from "@/pages/notes/components/MobileNotesShell/components/MobileDrawer/components/DrawerCloseButton/DrawerCloseButton.module.css"

interface DrawerCloseButtonProps {
    onClose: () => void
}

export default function DrawerCloseButton({onClose}: DrawerCloseButtonProps) {
    return (
        <div className={cls.DrawerCloseButtonContainer}>
            <Button
                variant="ghost"
                className={cls.CloseBtn}
                onClick={onClose}
                aria-label="Close sidebar">
                <ChevronLeftIcon/>
            </Button>
        </div>
    )
}
