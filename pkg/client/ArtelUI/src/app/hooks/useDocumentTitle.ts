import {useEffect} from "react"

export function useDocumentTitle(title: string) {
    useEffect(() => {
        document.title = `${title} — Artel`
        return () => {
            document.title = "Artel"
        }
    }, [title])
}
