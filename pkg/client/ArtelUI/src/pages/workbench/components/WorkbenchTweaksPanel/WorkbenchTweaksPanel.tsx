import {useEffect} from "react"

import {cn} from "@/app/utils/cn.ts"
import CloseIcon from "@/icons/common/CloseIcon.tsx"
import IconToggleButton
    from "@/pages/workbench/components/WorkbenchTopbar/components/IconToggleButton/IconToggleButton.tsx"
import type {WorkbenchMode} from "@/pages/workbench/processes/useWorkbenchModeControls.ts"
import type {WorkbenchContext} from "@/pages/workbench/processes/workbenchContext.ts"
import TweaksSystemPromptSection
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/workbench/components/WorkbenchTweaksPanel/components/TweaksSystemPromptSection/TweaksSystemPromptSection.tsx"
import TweaksThemeSection
    from "@/pages/workbench/components/WorkbenchTweaksPanel/components/TweaksThemeSection/TweaksThemeSection.tsx"
import TweaksMaxTokensSection
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/workbench/components/WorkbenchTweaksPanel/components/TweaksMaxTokensSection/TweaksMaxTokensSection.tsx"
import TweaksContextSection
    from "@/pages/workbench/components/WorkbenchTweaksPanel/components/TweaksContextSection/TweaksContextSection.tsx"
import TweaksConnectionsSection
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/workbench/components/WorkbenchTweaksPanel/components/TweaksConnectionsSection/TweaksConnectionsSection.tsx"
import cls from "@/pages/workbench/components/WorkbenchTweaksPanel/WorkbenchTweaksPanel.module.css"

interface Props {
    ctx: WorkbenchContext
    effectiveMode: WorkbenchMode | "picking"
    status: string
    vaultId?: string
}

// Right-side, in-flow, width-animated Tweaks panel. Always mounted (width: 0 when
// closed) so the open/close width transition runs; the section bodies mount only
// while open, same as TractCanvasInspector. Being the last flex child of
// .MainColumn, it paints over the chat content with no z-index. Escape closes it.
//
// >6 props isn't hit, but Props stays an object for parity with the other
// workbench panels.
export default function WorkbenchTweaksPanel(props: Props) {
    const {ctx, effectiveMode, status, vaultId} = props

    const {tweaksOpen, closeTweaks} = ctx
    useEffect(() => {
        if (!tweaksOpen) return
        function onKeyDown(e: KeyboardEvent) {
            if (e.key === "Escape") closeTweaks()
        }
        window.addEventListener("keydown", onKeyDown)
        return () => window.removeEventListener("keydown", onKeyDown)
    }, [tweaksOpen, closeTweaks])

    return (
        <div className={cn(cls.WorkbenchTweaksPanelContainer, cls.Panel, ctx.tweaksOpen && cls.Open)}>
            <div className={cls.Head}>
                <h2>Tweaks</h2>
                <IconToggleButton icon={<CloseIcon/>} label="Close tweaks" onClick={ctx.closeTweaks}/>
            </div>
            <div className={cls.Body}>
                {ctx.tweaksOpen && effectiveMode === "api" && (
                    <TweaksSystemPromptSection vaultId={vaultId}/>
                )}
                {ctx.tweaksOpen && <TweaksThemeSection/>}
                {ctx.tweaksOpen && <TweaksMaxTokensSection/>}
                {ctx.tweaksOpen && <TweaksContextSection/>}
                {ctx.tweaksOpen && (
                    <TweaksConnectionsSection effectiveMode={effectiveMode} status={status}/>
                )}
            </div>
        </div>
    )
}
