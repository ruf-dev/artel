import type {ComponentType} from "react"
import type {DropdownOption} from "@vervstack/chures"
import {Button, Dropdown, Input} from "@vervstack/chures"
import cls from "@/dialogs/AddTriggerDialog/AddTriggerDialog.module.css"

import {cn} from "@/app/utils/cn.ts"
import {useTracts, useTriggerSources} from "@/app/hooks/Tracts.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {providerEnumFor} from "@/components/TriggerPanel/triggerLabels.ts"

import TokenRevealDialog from "@/dialogs/TokenRevealDialog/TokenRevealDialog.tsx"
import {fieldsToSchemaNode, useAddTriggerDialog} from "@/dialogs/AddTriggerDialog/addTriggerDialogContext.ts"
import DialogHeaderWithClose from "@/dialogs/AddTriggerDialog/components/DialogHeaderWithClose.tsx"
import WebhookPicker from "@/components/WebhookPicker/WebhookPicker.tsx"
import SchemaBuilder from "@/dialogs/AddTriggerDialog/screens/SchemaBuilder.tsx"

const KIND_OPTIONS: DropdownOption[] = [
    {id: "webhook", name: "Webhook"},
    {id: "manual", name: "Manual"},
]

type KindSettingsProps = {
    source: string
    onSourceChange: (source: string) => void
}

const SettingsByKind: Record<"webhook" | "manual", ComponentType<KindSettingsProps> | undefined> = {
    webhook: WebhookPicker,
    manual: undefined,
}

export default function CreateScreen() {
    const {CloseDialog, OpenDialog} = useDialog()
    const {triggerSources} = useTriggerSources()

    const {connections} = useExternalConnections()
    const bakeError = useBakeError()

    const context = useAddTriggerDialog()

    function pickKind(opt: string[]) {
        if (opt[0] == 'webhook' || opt[0] == 'manual') {
            context.setKind(opt[0])
        }
    }

    const KindSettings = SettingsByKind[context.kind];

    const selectedSource = context.kind === "webhook" ? triggerSources.find(s => s.key === context.source) : undefined
    const hasProviderConnection = selectedSource?.provider
        ? connections.some(c => c.provider === providerEnumFor(selectedSource.provider))
        : true

    function handleCreate() {
        if (!context.name.trim()) return

        context.setSaving(true)

        const source = context.kind === "webhook" ? context.source : "generic"

        const payloadSchema = context.kind === "webhook" && context.source !== "generic"
            ? selectedSource?.payloadSchema ?? {properties: {}}
            : fieldsToSchemaNode(context.fields)

        useTracts
            .getState()
            .createAndLinkTrigger(context.name.trim(), context.kind, source, payloadSchema, context.tractUuid)
            .then(result => {
                CloseDialog()
                if (context.kind === "webhook" && result.webhookToken) {
                    setTimeout(() => OpenDialog(<TokenRevealDialog webhookUrl={result.webhookUrl}
                                                                   webhookToken={result.webhookToken}/>), 0)
                }
            })
            .catch(err => bakeError("Failed to create trigger", err))
            .finally(() => context.setSaving(false))
    }

    return (
        <div className={cls.DialogContainer} role="dialog" aria-modal="true">
            <DialogHeaderWithClose title="New trigger"/>

            <div className={cls.Body}>
                <div className={cls.Button}>
                    <Input
                        label="Name"
                        value={context.name}
                        setValue={context.setName}/>
                </div>

                <Dropdown
                    label={'Kind'}
                    options={KIND_OPTIONS}
                    onChange={pickKind}
                    value={[context.kind]}
                />

                {KindSettings && (
                    <KindSettings
                        source={context.source}
                        onSourceChange={context.setSource}/>
                )}

                {(context.kind === "manual" || context.source) && (
                    <SchemaBuilder
                        schemaId={context.kind === "webhook" && context.source !== "generic" ? context.source : ""}
                        fields={context.fields}
                        onChange={context.setFields}/>
                )}
            </div>

            <div className={cn(cls.DialogActions, cls.DialogActionsEnd)}>
                <Button
                    variant="ghost"
                    onClick={() => context.setStep("mode")}
                    disabled={context.saving}>
                    Back
                </Button>
                <Button
                    variant="primary"
                    onClick={handleCreate}
                    disabled={context.saving || !context.name.trim() || !hasProviderConnection}>
                    {context.saving ? "Creating…" : "Create & link"}
                </Button>
            </div>
        </div>
    )
}


