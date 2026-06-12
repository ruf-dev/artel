import { useEffect, useRef, useState } from "react"
import cls from "./NoteEditor.module.css"
import LineNumbers from "@/pages/notes/components/LineNumbers/LineNumbers.tsx"

interface NoteEditorProps {
    content: string
    onChange: (content: string) => void
    scrollTopRef: React.MutableRefObject<number>
}

export default function NoteEditor({ content, onChange, scrollTopRef }: NoteEditorProps) {
    const textareaRef = useRef<HTMLTextAreaElement>(null)
    const [scrollTop, setScrollTop] = useState(0)

    const lineCount = Math.max(content.split('\n').length, 20)

    useEffect(() => {
        if (textareaRef.current) {
            textareaRef.current.scrollTop = scrollTopRef.current
            setScrollTop(scrollTopRef.current)
        }
    }, [])

    function handleScroll() {
        const el = textareaRef.current
        if (!el) return
        scrollTopRef.current = el.scrollTop
        setScrollTop(el.scrollTop)
    }

    function handleKeyDown(e: React.KeyboardEvent<HTMLTextAreaElement>) {
        const el = e.currentTarget
        const meta = e.metaKey || e.ctrlKey

        if (e.key === 'Tab') {
            e.preventDefault()
            const { selectionStart, selectionEnd, value } = el
            if (!e.shiftKey) {
                const next = value.slice(0, selectionStart) + '  ' + value.slice(selectionEnd)
                onChange(next)
                requestAnimationFrame(() => {
                    el.selectionStart = selectionStart + 2
                    el.selectionEnd = selectionStart + 2
                })
            } else {
                const lineStart = value.lastIndexOf('\n', selectionStart - 1) + 1
                const leading = value.slice(lineStart, lineStart + 2)
                const spaces = leading.replace(/[^ ].*/, '').length
                if (spaces > 0) {
                    const remove = Math.min(spaces, 2)
                    const next = value.slice(0, lineStart) + value.slice(lineStart + remove)
                    onChange(next)
                    requestAnimationFrame(() => {
                        el.selectionStart = Math.max(selectionStart - remove, lineStart)
                        el.selectionEnd = Math.max(selectionEnd - remove, lineStart)
                    })
                }
            }
            return
        }

        if (meta && e.key === 'b') {
            e.preventDefault()
            wrapSelection(el, '**')
            return
        }

        if (meta && e.key === 'i') {
            e.preventDefault()
            wrapSelection(el, '*')
            return
        }
    }

    function wrapSelection(el: HTMLTextAreaElement, wrap: string) {
        const { selectionStart, selectionEnd, value } = el
        const selected = value.slice(selectionStart, selectionEnd)
        let replacement: string
        let nextStart: number
        let nextEnd: number

        if (selected.startsWith(wrap) && selected.endsWith(wrap)) {
            replacement = selected.slice(wrap.length, -wrap.length)
            nextStart = selectionStart
            nextEnd = selectionStart + replacement.length
        } else {
            replacement = `${wrap}${selected || 'text'}${wrap}`
            nextStart = selectionStart
            nextEnd = selectionStart + replacement.length
        }

        const next = value.slice(0, selectionStart) + replacement + value.slice(selectionEnd)
        onChange(next)
        requestAnimationFrame(() => {
            el.selectionStart = nextStart
            el.selectionEnd = nextEnd
        })
    }

    return (
        <div className={cls.NoteEditorContainer}>
            <div className={cls.IconSidebar} />
            <LineNumbers lineCount={lineCount} scrollTop={scrollTop} />
            <div className={cls.EditorArea}>
                <textarea
                    ref={textareaRef}
                    className={cls.Textarea}
                    value={content}
                    onChange={e => onChange(e.target.value)}
                    onScroll={handleScroll}
                    onKeyDown={handleKeyDown}
                    placeholder="Start writing…"
                    spellCheck={false}
                />
            </div>
        </div>
    )
}
