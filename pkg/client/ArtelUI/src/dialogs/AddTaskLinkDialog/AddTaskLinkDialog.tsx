import {useState} from "react"
import {Button, Dropdown, ModalClose} from "@vervstack/chures"

import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import Input from "@/components/atoms/Input/Input.tsx"
import cls from "@/dialogs/AddTaskLinkDialog/AddTaskLinkDialog.module.css"

export interface RoadmapLinkTarget {
    id: string
    name: string
    url: string
}

// Kept local (not imported from pages/roadmap/processes/taskLinkParser.ts) — dialogs must not
// import from pages (see eslint import-x/no-restricted-paths), and this dialog only ever writes
// 3 of the 4 relations a comment can express (the reader also understands "Blocked by", written
// by hand or by the other card's own "Blocks" link).
type WritableRelation = "depends_on" | "blocks" | "related_to"

const RELATION_OPTIONS: {id: WritableRelation; name: string}[] = [
    {id: "depends_on", name: "Depends on"},
    {id: "blocks", name: "Blocks"},
    {id: "related_to", name: "Related to"},
]

const RELATION_LABEL: Record<WritableRelation, string> = {
    depends_on: "Depends on",
    blocks: "Blocks",
    related_to: "Related to",
}

interface Props {
    trackerId: string
    sourceCardId: string
    /** Nodes already loaded in the caller's graph, mapped down to just what this dialog needs —
     * it doesn't know (or need) the full RoadmapNode shape. */
    targets: RoadmapLinkTarget[]
    onLinked: () => void
}

export default function AddTaskLinkDialog({trackerId, sourceCardId, targets, onLinked}: Props) {
    const {CloseDialog} = useDialog()
    const bakeError = useBakeError()
    const {executeMomTool} = useMcpKeys()

    const [relation, setRelation] = useState<WritableRelation>("depends_on")
    const [targetNodeId, setTargetNodeId] = useState("")
    const [url, setUrl] = useState("")
    const [saving, setSaving] = useState(false)

    const pickedTarget = targets.find(t => t.id === targetNodeId)
    const effectiveUrl = pickedTarget ? pickedTarget.url : url.trim()
    const canSubmit = effectiveUrl.length > 0 && !saving

    function handleSubmit() {
        if (!canSubmit) return
        setSaving(true)
        const text = `${RELATION_LABEL[relation]}: ${effectiveUrl}`
        executeMomTool("trello", "add_comment", trackerId, {card_id: sourceCardId, text})
            .then(() => {
                onLinked()
                CloseDialog()
            })
            .catch(err => bakeError("Failed to add link", err))
            .finally(() => setSaving(false))
    }

    return (
        <div
            className={cls.AddTaskLinkDialogContainer}
            role="dialog"
            aria-modal="true"
            aria-labelledby="addTaskLinkTitle"
        >
            <div className={cls.Head}>
                <h2 className={cls.Title} id="addTaskLinkTitle">Add task link</h2>
                <ModalClose onClick={CloseDialog} disabled={saving}/>
            </div>

            <label className={cls.Field}>
                <span className={cls.FieldLabel}>Relation</span>
                <Dropdown
                    options={RELATION_OPTIONS}
                    value={[relation]}
                    onChange={v => setRelation((v[0] as WritableRelation) || "depends_on")}
                />
            </label>

            <label className={cls.Field}>
                <span className={cls.FieldLabel}>Target — pick a card already on the canvas</span>
                <Dropdown
                    placeholder="Choose a loaded card…"
                    options={targets.map(t => ({id: t.id, name: t.name}))}
                    value={targetNodeId ? [targetNodeId] : []}
                    onChange={v => {
                        setTargetNodeId(v[0] ?? "")
                        setUrl("")
                    }}
                    emptyHint="No cards loaded yet — paste a URL below instead."
                />
            </label>

            <div className={cls.Divider}>or</div>

            <label className={cls.Field}>
                <span className={cls.FieldLabel}>Paste a Trello card URL</span>
                <Input
                    placeholder="https://trello.com/c/…"
                    value={url}
                    setValue={v => {
                        setUrl(v)
                        setTargetNodeId("")
                    }}
                    disabled={saving}
                    autoComplete="off"
                />
            </label>

            <div className={cls.Actions}>
                <Button variant="secondary" onClick={CloseDialog} disabled={saving}>Cancel</Button>
                <Button variant="primary" onClick={handleSubmit} disabled={!canSubmit}>
                    {saving ? "Adding…" : "Add link"}
                </Button>
            </div>
        </div>
    )
}
