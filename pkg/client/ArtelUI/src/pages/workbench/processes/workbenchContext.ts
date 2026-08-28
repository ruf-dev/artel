import {useCallback, useState} from "react"

// Which Tweaks-panel section a caller wants brought into view when it opens the
// panel. The panel scrolls/opens to the matching section; `undefined` just opens
// it at the top. Threaded down as one `WorkbenchContext` object (not a store, not
// N callbacks) from WorkbenchPage into the topbar and the composer chip row.
export type TweaksSection = "system" | "theme" | "tokens" | "context" | "connections"

export interface WorkbenchContext {
    tweaksOpen: boolean
    tweaksSection?: TweaksSection
    openTweaks: (section?: TweaksSection) => void
    closeTweaks: () => void
}

// Owns the Tweaks-panel open/section state for the workbench page. `openTweaks`
// opens the panel and records the requested section; `closeTweaks` just hides it
// and leaves the last section recorded so a re-open lands where it was.
export function useWorkbenchContext(): WorkbenchContext {
    const [tweaksOpen, setTweaksOpen] = useState(false)
    const [tweaksSection, setTweaksSection] = useState<TweaksSection | undefined>(undefined)

    const openTweaks = useCallback((section?: TweaksSection) => {
        setTweaksSection(section)
        setTweaksOpen(true)
    }, [])

    const closeTweaks = useCallback(() => {
        setTweaksOpen(false)
    }, [])

    return {tweaksOpen, tweaksSection, openTweaks, closeTweaks}
}
