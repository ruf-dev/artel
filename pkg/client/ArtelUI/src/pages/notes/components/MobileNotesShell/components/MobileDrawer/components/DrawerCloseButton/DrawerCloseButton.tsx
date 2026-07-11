import {Button} from "@vervstack/chures"

import ChevronLeftIcon from "@/pages/notes/components/icons/ChevronLeftIcon.tsx"
import cls
// eslint-disable-next-line max-len -- path too long to fit even unindented, can't shorten without changing import
from "@/pages/notes/components/MobileNotesShell/components/MobileDrawer/components/DrawerCloseButton/DrawerCloseButton.module.css"

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
