import {Button} from "@vervstack/chures"
import {MorphIcon} from "morphicons/react"
import {Settings, X} from "lucide"

import cls from "@/pages/workbench/components/WorkbenchTweaksPanel/components/SettingsFab/SettingsFab.module.css"

interface Props {
    open: boolean
    onToggle: () => void
}

// Floating bottom-right button that opens/closes the Tweaks overlay. A Settings
// glyph morphs to an X while the panel is open. Fixed-positioned and painted last
// (the panel portals to document.body), so it needs no z-index.
export default function SettingsFab({open, onToggle}: Props) {
    return (
        <div className={cls.SettingsFabContainer}>
            <Button
                variant="secondary"
                className={cls.Btn}
                aria-expanded={open}
                aria-label={open ? "Close settings" : "Settings"}
                onClick={onToggle}
            >
                <MorphIcon icon={open ? X : Settings} size={20} strokeWidth={1.6} className={cls.Icon}/>
            </Button>
        </div>
    )
}
