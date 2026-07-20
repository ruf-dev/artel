import {useState} from "react"
import {Button, Input} from "@vervstack/chures"

import cls from "@/pages/tract-templates/segments/ContentSegment/ContentSegment.module.css"
import {useDialog} from "@/app/hooks/Dialog"
import {useTractTemplates} from "@/app/hooks/TractTemplates.ts"
import TractTemplateCard from "@/pages/tract-templates/widgets/TractTemplateCard/TractTemplateCard.tsx"
import InstantiateTemplateDialog from "@/dialogs/InstantiateTemplateDialog/InstantiateTemplateDialog.tsx"

export default function ContentSegment() {
    const {templates, loading, mineOnly, category, setMineOnly, setCategory} = useTractTemplates()
    const {OpenDialog} = useDialog()
    const [categoryDraft, setCategoryDraft] = useState(category)

    function commitCategory() {
        if (categoryDraft.trim() !== category) {
            setCategory(categoryDraft.trim())
        }
    }

    return (
        <div className={cls.ContentSegmentContainer}>
            <div className={cls.Filters}>
                <div className={cls.Toggle}>
                    <Button
                        variant={mineOnly ? "ghost" : "primary"}
                        onClick={() => setMineOnly(false)}
                    >
                        All
                    </Button>
                    <Button
                        variant={mineOnly ? "primary" : "ghost"}
                        onClick={() => setMineOnly(true)}
                    >
                        Mine
                    </Button>
                </div>
                <Input
                    className={cls.CategoryInput}
                    value={categoryDraft}
                    setValue={setCategoryDraft}
                    onBlur={commitCategory}
                    onKeyDown={e => {
                        if (e.key === "Enter") commitCategory()
                    }}
                    placeholder="Filter by category…"
                />
            </div>

            {loading ? (
                <p className={cls.Empty}>Loading…</p>
            ) : templates.length === 0 ? (
                <div className={cls.EmptyState}>
                    <h2 className={cls.EmptyStateTitle}>No templates yet</h2>
                    <p className={cls.EmptyStateSub}>
                        {mineOnly
                            ? "You haven't published any tracts as templates yet."
                            : "No public templates match this filter."}
                    </p>
                </div>
            ) : (
                <div className={cls.Grid}>
                    {templates.map(t => (
                        <TractTemplateCard
                            key={t.uuid}
                            template={t}
                            mine={mineOnly}
                            onUse={() => OpenDialog(<InstantiateTemplateDialog templateUuid={t.uuid}/>)}
                        />
                    ))}
                </div>
            )}
        </div>
    )
}
