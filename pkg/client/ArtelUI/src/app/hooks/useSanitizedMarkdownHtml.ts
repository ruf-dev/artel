import {useMemo} from "react"
import {marked} from "marked"
import DOMPurify from "dompurify"

export function useSanitizedMarkdownHtml(content: string | null): string | null {
    return useMemo(() => {
        if (content === null) return null
        if (content === '') return ''
        const rawHtml = marked.parse(content) as string
        return DOMPurify.sanitize(rawHtml)
    }, [content])
}
