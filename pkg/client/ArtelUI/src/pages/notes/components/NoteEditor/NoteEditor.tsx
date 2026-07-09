import { useEffect, useRef, useState } from "react"

import LineNumbers from "@/pages/notes/components/LineNumbers/LineNumbers.tsx"
import BoldIcon from "@/pages/notes/components/icons/BoldIcon.tsx"
import ItalicIcon from "@/pages/notes/components/icons/ItalicIcon.tsx"
import LinkIcon from "@/pages/notes/components/icons/LinkIcon.tsx"
import HeadingIcon from "@/pages/notes/components/icons/HeadingIcon.tsx"
import CodeIcon from "@/pages/notes/components/icons/CodeIcon.tsx"
import { handleEditorKeyDown } from "@/pages/notes/components/NoteEditor/processes/noteEditorKeyHandlers.ts"
import cls from "@/pages/notes/components/NoteEditor/NoteEditor.module.css"

interface NoteEditorProps {
    content: string
    onChange: (content: string) => void
    scrollTopRef: React.MutableRefObject<number>
    fontScale?: number
    onEscape?: () => void
}

export default function NoteEditor({ content, onChange, scrollTopRef, fontScale = 1, onEscape }: NoteEditorProps) {
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

    return (
        <div className={cls.NoteEditorContainer}>
            <div className={cls.IconBar}>
                <div className={cls.IconBarButton} data-tooltip-id="root-tooltip" data-tooltip-content="Still working on it"><BoldIcon /></div>
                <div className={cls.IconBarButton} data-tooltip-id="root-tooltip" data-tooltip-content="Still working on it"><ItalicIcon /></div>
                <div className={cls.IconBarButton} data-tooltip-id="root-tooltip" data-tooltip-content="Still working on it"><HeadingIcon /></div>
                <div className={cls.IconBarButton} data-tooltip-id="root-tooltip" data-tooltip-content="Still working on it"><LinkIcon /></div>
                <div className={cls.IconBarButton} data-tooltip-id="root-tooltip" data-tooltip-content="Still working on it"><CodeIcon /></div>
            </div>
            <div className={cls.EditorRow}>
                <LineNumbers lineCount={lineCount} scrollTop={scrollTop} lineHeight={21 * fontScale} />
                <div className={cls.EditorArea}>
                    <textarea
                        ref={textareaRef}
                        className={cls.Textarea}
                        style={{ fontSize: 12 * fontScale, lineHeight: 1.75 }}
                        value={content}
                        onChange={e => onChange(e.target.value)}
                        onScroll={handleScroll}
                        onKeyDown={e => handleEditorKeyDown(e, onChange, onEscape)}
                        placeholder="Start writing…"
                        spellCheck={false}
                    />
                </div>
            </div>
        </div>
    )
}
