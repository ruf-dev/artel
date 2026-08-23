import {useState} from "react"

interface Params {
    selectTab: (tabId: string) => Promise<unknown>
    createTab: () => Promise<unknown>
    closeTab: (tabId: string) => Promise<unknown>
    bakeError: (title: string, err: unknown) => void
}

// Bundles the terminal-tab action handlers (wrapped with error toasting) and the
// history-sidebar open/close state for WorkbenchPage — split out purely to keep
// WorkbenchPage.tsx's render function under the max-lines-per-function lint limit.
export function useWorkbenchPanelControls({selectTab, createTab, closeTab, bakeError}: Params) {
    const [historyOpen, setHistoryOpen] = useState(false)

    function handleSelectTab(tabId: string) {
        selectTab(tabId).catch(e => bakeError("Failed to select tab", e))
    }

    function handleCreateTab() {
        createTab().catch(e => bakeError("Failed to create tab", e))
    }

    function handleCloseTab(tabId: string) {
        closeTab(tabId).catch(e => bakeError("Failed to close tab", e))
    }

    function toggleHistory() {
        setHistoryOpen(open => !open)
    }

    function closeHistory() {
        setHistoryOpen(false)
    }

    return {historyOpen, handleSelectTab, handleCreateTab, handleCloseTab, toggleHistory, closeHistory}
}
