import {useEffect, useState} from "react"
import {useNavigate} from "react-router-dom"
import {Button, Input, LoadingWrapper} from "@vervstack/chures"

import cls from "@/dialogs/InstantiateTemplateDialog/InstantiateTemplateDialog.module.css"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import {TractTemplate, tractsService} from "@/processes/Tracts.ts"
import {useTemplateConnections} from "@/dialogs/InstantiateTemplateDialog/hooks/useTemplateConnections.ts"
import ConnectionSection from "@/dialogs/InstantiateTemplateDialog/components/ConnectionSection/ConnectionSection.tsx"

interface Props {
    templateUuid: string
}

export default function InstantiateTemplateDialog({templateUuid}: Props) {
    const {CloseDialog} = useDialog()
    const bakeError = useBakeError()
    const navigate = useNavigate()

    const [template, setTemplate] = useState<TractTemplate | null>(null)
    const [loading, setLoading] = useState(true)
    const [name, setName] = useState("")
    const [description, setDescription] = useState("")
    const [submitting, setSubmitting] = useState(false)

    const {requirements, candidatesByKey, connectionsByKey, setConnectionsByKey} = useTemplateConnections(
        template?.definition.steps ?? [],
    )

    useEffect(() => {
        setLoading(true)
        tractsService.getTemplate(templateUuid)
            .then(t => {
                setTemplate(t)
                setName(t.name)
                setDescription(t.description)
            })
            .catch(err => bakeError("Failed to load template", err))
            .finally(() => setLoading(false))
    }, [templateUuid, bakeError])

    const canSubmit = !!name.trim() && requirements.every(req => !!connectionsByKey[req.key])

    function handleSubmit() {
        setSubmitting(true)
        tractsService.instantiateTemplate(templateUuid, name.trim(), description, connectionsByKey)
            .then(({tract}) => {
                CloseDialog()
                navigate(`/tracts/${tract.uuid}`)
            })
            .catch(err => bakeError("Failed to instantiate template", err))
            .finally(() => setSubmitting(false))
    }

    return (
        <div className={cls.InstantiateTemplateDialogContainer} role="dialog" aria-modal="true">
            <LoadingWrapper isLoading={loading}>
                <h2 className={cls.DialogTitle}>Use template</h2>
                <Input
                    className={cls.DialogInputWrapper}
                    inputClassName={cls.DialogInput}
                    value={name}
                    setValue={setName}
                    placeholder="Tract name"
                    autoFocus
                    disabled={submitting}
                />
                <Input
                    className={cls.DialogInputWrapper}
                    inputClassName={cls.DialogInput}
                    value={description}
                    setValue={setDescription}
                    placeholder="Description (optional)"
                    disabled={submitting}
                />
                {requirements.length > 0 && (
                    <p className={cls.Requires}>
                        Requires <span className={cls.MomNames}>{requirements.map(r => r.key).join(", ")}</span>
                        {requirements.length > 1 ? " connections" : " connection"}
                    </p>
                )}
                {requirements.length > 0 && (
                    <div className={cls.Connections}>
                        {requirements.map(req => (
                            <ConnectionSection
                                key={req.key}
                                mcp={req.key}
                                connections={candidatesByKey[req.key] ?? []}
                                selectedId={connectionsByKey[req.key] ?? ""}
                                onSelect={id => setConnectionsByKey(prev => ({...prev, [req.key]: id}))}
                            />
                        ))}
                    </div>
                )}
                <div className={cls.DialogActions}>
                    <Button variant="ghost" onClick={CloseDialog} disabled={submitting}>Cancel</Button>
                    <Button variant="primary" onClick={handleSubmit} disabled={submitting || !canSubmit}>
                        {submitting ? "…" : "Create tract"}
                    </Button>
                </div>
            </LoadingWrapper>
        </div>
    )
}
