import {Button} from "@vervstack/chures"
import {MorphIcon} from "morphicons/react"
import {Menu, X} from "lucide"

// eslint-disable-next-line max-len -- import path too long to wrap under 120 chars
import cls from "@/pages/workbench/components/WorkbenchSidebar/components/SidebarToggleButton/SidebarToggleButton.module.css"

interface Props {
    open: boolean
    onToggle: () => void
}

// Collapse/expand the workbench sidebar. Lives at the top-right of WorkbenchSidebar
// so it stays visible when the sidebar column collapses to the narrow rail. A
// Menu glyph morphs to an X once the sidebar is open.
export default function SidebarToggleButton({open, onToggle}: Props) {
    return (
        <div className={cls.SidebarToggleButtonContainer}>
            <Button
                variant="unstyled"
                className={cls.Btn}
                aria-expanded={open}
                aria-label="Toggle conversations"
                onClick={onToggle}
            >
                <MorphIcon icon={open ? X : Menu} size={20} strokeWidth={1.6} className={cls.Icon}/>
            </Button>
        </div>
    )
}
