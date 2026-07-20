import {useEffect, useMemo, useState} from "react"
import {useNavigate} from "react-router-dom"
import {Button, Input, LoadingWrapper} from "@vervstack/chures"

import cls from "@/dialogs/InstantiateTemplateDialog/InstantiateTemplateDialog.module.css"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import {TractTemplate, tractsService} from "@/processes/Tracts.ts"
import {requiredConnections} from "@/dialogs/InstantiateTemplateDialog/processes/requiredConnections.ts"
import ConnectionSection from "@/dialogs/InstantiateTemplateDialog/components/ConnectionSection/ConnectionSection.tsx"

interface Props {
    templateUuid: string
}

export default function InstantiateTemplateDialog({templateUuid}: Props) {
    const {CloseDialog} = useDialog()
    const bakeError = useBakeError()
    const navigate = useNavigate()
    const {momCandidates, fetchMomCandidates} = useMcpKeys()

    const [template, setTemplate] = useState<TractTemplate | null>(null)
    const [loading, setLoading] = useState(true)
    const [name, setName] = useState("")
    const [description, setDescription] = useState("")
    const [connectionsByMom, setConnectionsByMom] = useState<Record<string, string>>({})
    const [submitting, setSubmitting] = useState(false)

    useEffect(() => {
        void fetchMomCandidates()
    }, [fetchMomCandidates])

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

    const requiredMoms = useMemo(
        () => template ? requiredConnections(template.definition.steps) : [],
        [template],
    )

    const canSubmit = !!name.trim() && requiredMoms.every(mcp => !!connectionsByMom[mcp])

    function handleSubmit() {
        setSubmitting(true)
        tractsService.instantiateTemplate(templateUuid, name.trim(), description, connectionsByMom)
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
                {requiredMoms.length > 0 && (
                    <p className={cls.Requires}>
                        Requires <span className={cls.MomNames}>{requiredMoms.join(", ")}</span>
                        {requiredMoms.length > 1 ? " connections" : " connection"}
                    </p>
                )}
                {requiredMoms.length > 0 && (
                    <div className={cls.Connections}>
                        {requiredMoms.map(mcp => (
                            <ConnectionSection
                                key={mcp}
                                mcp={mcp}
                                connections={momCandidates.find(c => c.name === mcp)?.connections ?? []}
                                selectedId={connectionsByMom[mcp] ?? ""}
                                onSelect={id => setConnectionsByMom(prev => ({...prev, [mcp]: id}))}
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
