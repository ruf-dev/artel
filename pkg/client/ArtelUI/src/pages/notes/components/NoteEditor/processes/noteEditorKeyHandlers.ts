export function wrapSelection(el: HTMLTextAreaElement, wrap: string, onChange: (content: string) => void) {
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

export function handleEditorKeyDown(
    e: React.KeyboardEvent<HTMLTextAreaElement>,
    onChange: (content: string) => void,
    onEscape?: () => void,
) {
    const el = e.currentTarget
    const meta = e.metaKey || e.ctrlKey

    if (e.key === 'Escape' && onEscape) {
        e.preventDefault()
        onEscape()
        return
    }

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
        wrapSelection(el, '**', onChange)
        return
    }

    if (meta && e.key === 'i') {
        e.preventDefault()
        wrapSelection(el, '*', onChange)
        return
    }
}
