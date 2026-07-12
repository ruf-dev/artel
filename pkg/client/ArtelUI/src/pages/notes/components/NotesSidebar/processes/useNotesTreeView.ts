import { useEffect, useState } from "react"

const STORAGE_KEY = "artel.notes.searchTreeView"

export function useNotesTreeView(): [boolean, () => void] {
    const [treeView, setTreeView] = useState<boolean>(
        () => localStorage.getItem(STORAGE_KEY) === "true"
    )

    useEffect(() => {
        localStorage.setItem(STORAGE_KEY, String(treeView))
    }, [treeView])

    function toggleTreeView() {
        setTreeView(v => !v)
    }

    return [treeView, toggleTreeView]
}
