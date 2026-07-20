import {useMemo, useState} from "react"
import {Button, Input} from "@vervstack/chures"

import cls from "@/dialogs/BrowseTemplatesDialog/BrowseTemplatesDialog.module.css"
import {useDialog} from "@/app/hooks/Dialog"
import {TractTemplateSummary} from "@/processes/Tracts.ts"
import TemplateRow from "@/dialogs/BrowseTemplatesDialog/components/TemplateRow/TemplateRow.tsx"

interface Props {
    templates: TractTemplateSummary[]
    loading: boolean
    onSelect: (template: TractTemplateSummary) => void
}

export default function ListScreen({templates, loading, onSelect}: Props) {
    const {CloseDialog} = useDialog()
    const [search, setSearch] = useState("")

    const filtered = useMemo(() => {
        const q = search.trim().toLowerCase()
        if (!q) return templates
        return templates.filter(t => t.name.toLowerCase().includes(q) || t.description.toLowerCase().includes(q))
    }, [templates, search])

    return (
        <div className={cls.BrowseTemplatesDialogContainer} role="dialog" aria-modal="true">
            <h2 className={cls.DialogTitle}>Explore templates</h2>
            <Input
                inputClassName={cls.SearchInput}
                value={search}
                setValue={setSearch}
                placeholder="Search templates…"
                autoFocus
            />
            <div className={cls.List}>
                {loading ? (
                    <p className={cls.Empty}>Loading…</p>
                ) : filtered.length === 0 ? (
                    <p className={cls.Empty}>No templates found.</p>
                ) : (
                    filtered.map(t => <TemplateRow key={t.uuid} template={t} onClick={() => onSelect(t)}/>)
                )}
            </div>
            <div className={cls.DialogActions}>
                <Button variant="ghost" onClick={CloseDialog}>Close</Button>
            </div>
        </div>
    )
}
