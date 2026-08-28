import React from "react"

import {cn} from "@/app/utils/cn.ts"
import cls from "@/pages/workbench/components/WorkbenchTweaksPanel/components/TweaksSection/TweaksSection.module.css"

interface Props {
    label: string
    children: React.ReactNode
    disabled?: boolean
}

// One labelled block inside the Tweaks panel body — mirrors the design mock's
// `.tp-section` (a `<label>` above a single control). `disabled` dims the whole
// block for the not-yet-wired placeholder sections (Max tokens, Context window).
export default function TweaksSection({label, children, disabled}: Props) {
    return (
        <div className={cn(cls.TweaksSectionContainer, disabled && cls.Disabled)}>
            <label className={cls.Label}>{label}</label>
            {children}
        </div>
    )
}
