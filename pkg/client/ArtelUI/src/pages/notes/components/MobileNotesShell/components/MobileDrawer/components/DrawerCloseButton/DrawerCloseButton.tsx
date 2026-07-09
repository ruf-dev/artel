import {Button} from "@vervstack/chures"

import CloseIcon from "@/pages/notes/components/icons/CloseIcon.tsx"
import cls from "@/pages/notes/components/MobileNotesShell/components/MobileDrawer/components/DrawerCloseButton/DrawerCloseButton.module.css"

interface DrawerCloseButtonProps {
    onClose: () => void
}

export default function DrawerCloseButton({onClose}: DrawerCloseButtonProps) {
    return (
        <Button variant="ghost" className={cls.DrawerCloseBtn} onClick={onClose} aria-label="Close sidebar">
            <CloseIcon/>
        </Button>
    )
}
