import {useEffect, useState} from "react"
import {Button, LoadingWrapper} from "@vervstack/chures"

import cls from "@/dialogs/BrowseTemplatesDialog/BrowseTemplatesDialog.module.css"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import {TractTemplate, TractTemplateSummary, tractsService} from "@/processes/Tracts.ts"
import TractsIcon from "@/segments/Topbar/components/icons/TractsIcon.tsx"
import InstantiateTemplateDialog from "@/dialogs/InstantiateTemplateDialog/InstantiateTemplateDialog.tsx"

interface Props {
    template: TractTemplateSummary
    onBack: () => void
}

export default function DetailScreen({template, onBack}: Props) {
    const {OpenDialog} = useDialog()
    const bakeError = useBakeError()
    const [full, setFull] = useState<TractTemplate | null>(null)
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        setLoading(true)
        tractsService.getTemplate(template.uuid)
            .then(setFull)
            .catch(err => bakeError("Failed to load template", err))
            .finally(() => setLoading(false))
    }, [template.uuid, bakeError])

    return (
        <div className={cls.BrowseTemplatesDialogContainer} role="dialog" aria-modal="true">
            <div className={cls.DetailHeader}>
                <TractsIcon className={cls.DetailIcon}/>
                <div className={cls.DetailHeaderInfo}>
                    <h2 className={cls.DialogTitle}>{template.name}</h2>
                    {template.category && <span className={cls.Chip}>{template.category}</span>}
                </div>
            </div>
            <LoadingWrapper isLoading={loading}>
                {template.description && <p className={cls.DetailDesc}>{template.description}</p>}
                <span className={cls.InstallCount}>
                    {template.installCount} {template.installCount === 1 ? "install" : "installs"}
                </span>
                {full && full.definition.steps.length > 0 && (
                    <>
                        <span className={cls.StepsLabel}>Steps</span>
                        <div className={cls.StepsPreview}>
                            {full.definition.steps.map(step => (
                                <span key={step.id} className={cls.Chip}>{step.name || step.type}</span>
                            ))}
                        </div>
                    </>
                )}
            </LoadingWrapper>
            <div className={cls.DialogActions}>
                <Button variant="ghost" onClick={onBack}>Back</Button>
                <Button
                    variant="primary"
                    onClick={() => OpenDialog(<InstantiateTemplateDialog templateUuid={template.uuid}/>)}
                >
                    Use this template
                </Button>
            </div>
        </div>
    )
}
