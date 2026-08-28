import {useCallback, useState} from "react"

// Owns the "attached vault files" set for the workbench chat composer. Threaded
// down as one plain `WorkbenchAttachments` object (not a store, not N callbacks)
// from WorkbenchPage into WorkbenchPanels -> ChatPanelShell, mirroring the
// workbenchContext.ts pattern. Stage 4 also feeds this same object into the
// sidebar Vault pane so clicking a file there calls `toggle`.

export interface WorkbenchAttachments {
    paths: string[]
    toggle: (path: string) => void
    remove: (path: string) => void
    clear: () => void
}

export function useWorkbenchAttachments(): WorkbenchAttachments {
    const [paths, setPaths] = useState<string[]>([])

    const toggle = useCallback((path: string) => {
        setPaths(prev => prev.includes(path) ? prev.filter(p => p !== path) : [...prev, path])
    }, [])

    const remove = useCallback((path: string) => {
        setPaths(prev => prev.filter(p => p !== path))
    }, [])

    const clear = useCallback(() => {
        setPaths([])
    }, [])

    return {paths, toggle, remove, clear}
}
