import { useCallback, useEffect, useRef, useState } from "react"
import { SaveStatus } from "@/pages/notes/components/SaveStatusIndicator/SaveStatusIndicator.tsx"

const AUTOSAVE_DELAY = 500

interface UseAutosaveOptions {
    noteId: string | null
    vaultId: string | null
    content: string
    savedContent: string
    onSave: (content: string) => Promise<void>
    onSavedContentChange: (content: string) => void
}

interface UseAutosaveResult {
    saveStatus: SaveStatus
    saveError: string | undefined
    forceSave: () => void
    isDirty: boolean
}

export function useAutosave({
    noteId,
    vaultId,
    content,
    savedContent,
    onSave,
    onSavedContentChange,
}: UseAutosaveOptions): UseAutosaveResult {
    const [saveStatus, setSaveStatus] = useState<SaveStatus>('idle')
    const [saveError, setSaveError] = useState<string | undefined>()
    const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null)
    const savedStatusTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null)
    const isSavingRef = useRef(false)
    const pendingContentRef = useRef<string | null>(null)

    const isDirty = content !== savedContent

    const doSave = useCallback(async (contentToSave: string) => {
        if (!noteId || !vaultId) return
        if (isSavingRef.current) {
            pendingContentRef.current = contentToSave
            return
        }

        setSaveStatus('saving')
        setSaveError(undefined)
        isSavingRef.current = true

        try {
            await onSave(contentToSave)
            onSavedContentChange(contentToSave)
            setSaveStatus('saved')
            if (savedStatusTimerRef.current) clearTimeout(savedStatusTimerRef.current)
            savedStatusTimerRef.current = setTimeout(() => setSaveStatus('idle'), 3000)

            if (noteId) {
                localStorage.removeItem(`artel.notes.draft.${noteId}`)
            }
        } catch (err) {
            setSaveStatus('error')
            setSaveError(err instanceof Error ? err.message : String(err))
        } finally {
            isSavingRef.current = false
            if (pendingContentRef.current !== null) {
                const next = pendingContentRef.current
                pendingContentRef.current = null
                void doSave(next)
            }
        }
    }, [noteId, vaultId, onSave, onSavedContentChange])

    const forceSave = useCallback(() => {
        if (timerRef.current) clearTimeout(timerRef.current)
        if (content !== savedContent) {
            void doSave(content)
        }
    }, [content, savedContent, doSave])

    useEffect(() => {
        if (!isDirty) return

        if (timerRef.current) clearTimeout(timerRef.current)
        setSaveStatus('dirty')

        timerRef.current = setTimeout(() => {
            void doSave(content)
        }, AUTOSAVE_DELAY)

        if (noteId) {
            localStorage.setItem(`artel.notes.draft.${noteId}`, content)
        }

        return () => {
            if (timerRef.current) clearTimeout(timerRef.current)
        }
    }, [content, isDirty, noteId])

    useEffect(() => {
        return () => {
            if (savedStatusTimerRef.current) clearTimeout(savedStatusTimerRef.current)
        }
    }, [])

    useEffect(() => {
        function handleBeforeUnload(e: BeforeUnloadEvent) {
            if (saveStatus === 'dirty' || saveStatus === 'saving') {
                e.preventDefault()
            }
        }

        if (saveStatus === 'dirty' || saveStatus === 'saving') {
            window.addEventListener('beforeunload', handleBeforeUnload)
        }
        return () => window.removeEventListener('beforeunload', handleBeforeUnload)
    }, [saveStatus])

    return { saveStatus, saveError, forceSave, isDirty }
}
