import {Button, ConfirmDialog} from "@vervstack/chures"

import cls from "@/pages/tract-templates/widgets/TractTemplateCard/TractTemplateCard.module.css"
import {TractTemplateSummary} from "@/processes/Tracts.ts"
import {useTractTemplates} from "@/app/hooks/TractTemplates.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"

interface Props {
    template: TractTemplateSummary
    mine: boolean
    onUse: () => void
}

export default function TractTemplateCard({template, mine, onUse}: Props) {
    const {unpublish} = useTractTemplates()
    const bakeError = useBakeError()
    const {OpenDialog, CloseDialog} = useDialog()

    function handleUnpublish(e: React.MouseEvent) {
        e.stopPropagation()
        OpenDialog(
            <ConfirmDialog
                title="Unpublish template"
                message={`Unpublish "${template.name}"? It will no longer be visible to other users.`}
                confirmLabel="Unpublish"
                danger
                onClose={CloseDialog}
                onConfirm={() => unpublish(template.uuid).catch(err => bakeError("Failed to unpublish template", err))}
            />
        )
    }

    return (
        <div className={cls.Card}>
            <div className={cls.Info}>
                <span className={cls.Name}>{template.name}</span>
                {template.description && <span className={cls.Desc}>{template.description}</span>}
            </div>
            <div className={cls.Meta}>
                {template.category && <span className={cls.Chip}>{template.category}</span>}
                <span className={cls.InstallCount}>
                    {template.installCount} {template.installCount === 1 ? "install" : "installs"}
                </span>
                {mine && <span className={cls.OwnerChip}>Yours</span>}
            </div>
            <div className={cls.Actions}>
                {mine && (
                    <Button variant="danger" onClick={handleUnpublish}>
                        Unpublish
                    </Button>
                )}
                <Button variant="primary" onClick={onUse}>
                    Use this template
                </Button>
            </div>
        </div>
    )
}
