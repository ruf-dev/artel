import { useSearchParams } from "react-router-dom"

const SEARCH_PARAM = "q"

export function useNotesSearchQuery(): [string, (value: string) => void] {
    const [searchParams, setSearchParams] = useSearchParams()
    const searchQuery = searchParams.get(SEARCH_PARAM) ?? ""

    function setSearchQuery(value: string) {
        setSearchParams(prev => {
            const next = new URLSearchParams(prev)
            if (value) {
                next.set(SEARCH_PARAM, value)
            } else {
                next.delete(SEARCH_PARAM)
            }
            return next
        }, {replace: true})
    }

    return [searchQuery, setSearchQuery]
}
