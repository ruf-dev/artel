import {useEffect} from "react"
import {createPortal} from "react-dom"

import {cn} from "@/app/utils/cn.ts"
import CloseIcon from "@/icons/common/CloseIcon.tsx"
import IconToggleButton
    from "@/pages/workbench/components/WorkbenchTopbar/components/IconToggleButton/IconToggleButton.tsx"
import type {WorkbenchMode} from "@/pages/workbench/processes/useWorkbenchModeControls.ts"
import type {WorkbenchContext} from "@/pages/workbench/processes/workbenchContext.ts"
import TweaksSystemPromptSection
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/workbench/components/WorkbenchTweaksPanel/components/TweaksSystemPromptSection/TweaksSystemPromptSection.tsx"
import TweaksMaxTokensSection
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/workbench/components/WorkbenchTweaksPanel/components/TweaksMaxTokensSection/TweaksMaxTokensSection.tsx"
import TweaksContextSection
    from "@/pages/workbench/components/WorkbenchTweaksPanel/components/TweaksContextSection/TweaksContextSection.tsx"
import TweaksConnectionsSection
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/workbench/components/WorkbenchTweaksPanel/components/TweaksConnectionsSection/TweaksConnectionsSection.tsx"
import SettingsFab
    from "@/pages/workbench/components/WorkbenchTweaksPanel/components/SettingsFab/SettingsFab.tsx"
import cls from "@/pages/workbench/components/WorkbenchTweaksPanel/WorkbenchTweaksPanel.module.css"

interface Props {
    ctx: WorkbenchContext
    effectiveMode: WorkbenchMode | "picking"
    status: string
    vaultId?: string
}

// Full-viewport Tweaks overlay: a blurred backdrop over the whole page plus a
// right-side panel that slides in. Portaled to document.body (appended after
// #root) so it paints above everything with no z-index — same pattern as
// WorkbenchSettingsMenu. Always mounted so the slide-out transition runs; the
// section bodies mount only while open. Escape and a backdrop click close it. The
// floating SettingsFab (its own portal) toggles the overlay open/closed.
//
// >6 props isn't hit, but Props stays an object for parity with the other
// workbench panels.
export default function WorkbenchTweaksPanel(props: Props) {
    const {ctx, effectiveMode, status, vaultId} = props
    const {tweaksOpen, openTweaks, closeTweaks} = ctx

    useEffect(() => {
        if (!tweaksOpen) return
        function onKeyDown(e: KeyboardEvent) {
            if (e.key === "Escape") closeTweaks()
        }
        window.addEventListener("keydown", onKeyDown)
        return () => window.removeEventListener("keydown", onKeyDown)
    }, [tweaksOpen, closeTweaks])

    return (
        <>
            {createPortal(
                <div className={cn(cls.WorkbenchTweaksPanelContainer, tweaksOpen && cls.Open)}>
                    <div className={cls.Backdrop} onClick={closeTweaks}/>
                    <div className={cls.Panel}>
                        <div className={cls.Head}>
                            <h2>Tweaks</h2>
                            <IconToggleButton icon={<CloseIcon/>} label="Close tweaks" onClick={closeTweaks}/>
                        </div>
                        <div className={cls.Body}>
                            {tweaksOpen && effectiveMode === "api" && (
                                <TweaksSystemPromptSection vaultId={vaultId}/>
                            )}
                            {tweaksOpen && <TweaksMaxTokensSection/>}
                            {tweaksOpen && <TweaksContextSection/>}
                            {tweaksOpen && (
                                <TweaksConnectionsSection effectiveMode={effectiveMode} status={status}/>
                            )}
                        </div>
                    </div>
                </div>,
                document.body,
            )}
            {createPortal(
                <SettingsFab
                    open={tweaksOpen}
                    onToggle={() => tweaksOpen ? closeTweaks() : openTweaks()}
                />,
                document.body,
            )}
        </>
    )
}
