import {useState} from "react"
import {Button} from "@vervstack/chures"

import cls from "@/dialogs/ConnectGenericDialog/ConnectGenericDialog.module.css"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import CredentialRow from "@/dialogs/ConnectGenericDialog/components/CredentialRow/CredentialRow.tsx"

interface CredentialField {
    key: string
    value: string
}

export default function ConnectGenericDialog({name, onSuccess}: {name: string; onSuccess: () => void}) {
    const {CloseDialog} = useDialog()
    const bakeError = useBakeError()
    const {addGenericConnection} = useExternalConnections()

    const [fields, setFields] = useState<CredentialField[]>([{key: "", value: ""}])
    const [submitting, setSubmitting] = useState(false)

    const hasValidField = fields.some(f => f.key.trim().length > 0)

    function updateField(index: number, patch: Partial<CredentialField>) {
        setFields(prev => prev.map((f, i) => (i === index ? {...f, ...patch} : f)))
    }

    function addField() {
        setFields(prev => [...prev, {key: "", value: ""}])
    }

    function removeField(index: number) {
        setFields(prev => prev.filter((_, i) => i !== index))
    }

    function handleSubmit() {
        const credentials: Record<string, string> = {}
        fields.forEach(f => {
            if (f.key.trim()) credentials[f.key.trim()] = f.value
        })
        setSubmitting(true)
        addGenericConnection({provider: name, credentials})
            .then(() => {
                onSuccess()
                CloseDialog()
            })
            .catch(e => bakeError("Failed to connect", e))
            .finally(() => setSubmitting(false))
    }

    return (
        <div className={cls.ConnectGenericDialogContainer} role="dialog" aria-modal="true">
            <h2 className={cls.DialogTitle}>Connect to {name}</h2>
            <p className={cls.Notice}>
                Enter the credential field names and values that <b>{name}</b>&apos;s author told you to provide —
                there&apos;s no machine-readable schema to pre-fill this form from.
            </p>
            <div className={cls.FieldsList}>
                {fields.map((f, i) => (
                    <CredentialRow
                        key={i}
                        keyValue={f.key}
                        value={f.value}
                        onKeyChange={v => updateField(i, {key: v})}
                        onValueChange={v => updateField(i, {value: v})}
                        onRemove={() => removeField(i)}
                        disabled={submitting}
                        canRemove={fields.length > 1}
                    />
                ))}
            </div>
            <Button variant="ghost" onClick={addField} disabled={submitting}>+ Add field</Button>
            <div className={cls.DialogActions}>
                <Button variant="ghost" onClick={CloseDialog} disabled={submitting}>Cancel</Button>
                <Button variant="primary" onClick={handleSubmit} disabled={submitting || !hasValidField}>
                    {submitting ? "Connecting…" : "Connect"}
                </Button>
            </div>
        </div>
    )
}
