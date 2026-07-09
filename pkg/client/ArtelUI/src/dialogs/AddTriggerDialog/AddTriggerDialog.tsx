import {type ReactNode, useEffect, useState} from "react"

import {TractCondition} from "@/processes/Tracts.ts"
import {useTracts} from "@/app/hooks/Tracts.ts"
import {
    AddTriggerDialogContext,
    AddTriggerDialogState,
    AddTriggerStep,
    emptySchemaField,
    SchemaFieldRow
} from "@/dialogs/AddTriggerDialog/addTriggerDialogContext.ts"
import ModeScreen from "@/dialogs/AddTriggerDialog/screens/ModeScreen.tsx"
import CreateScreen from "@/dialogs/AddTriggerDialog/screens/CreateScreen.tsx"
import LinkScreen from "@/dialogs/AddTriggerDialog/screens/LinkScreen.tsx"

export default function AddTriggerDialog({tractUuid, linkedUuids}: { tractUuid: string; linkedUuids: Set<string> }) {
    const {fetchTriggers} = useTracts()

    const [step, setStep] = useState<AddTriggerStep>("mode")
    const [saving, setSaving] = useState(false)

    const [name, setName] = useState("")
    const [kind, setKind] = useState<"webhook" | "manual">("webhook")
    const [source, setSource] = useState("generic")
    const [fields, setFields] = useState<SchemaFieldRow[]>([emptySchemaField()])

    const [selectedTriggerId, setSelectedTriggerId] = useState("")
    const [filters, setFilters] = useState<TractCondition[]>([])

    // Refetch on every open — this dialog mounts fresh each time OpenDialog is called
    // (see app/hooks/Dialog.ts), so a plain mount effect refreshes stale triggers instead
    // of reusing whatever TriggerPanel's own mount fetch loaded at page-load time. Trigger
    // sources are covered by useTriggerSources' own react-query fetch-on-mount in CreateScreen.
    useEffect(() => {
        void fetchTriggers()
    }, [fetchTriggers])

    const value: AddTriggerDialogState = {
        tractUuid, linkedUuids,
        step, setStep,
        saving, setSaving,
        name, setName,
        kind, setKind,
        source, setSource,
        fields, setFields,
        selectedTriggerId, setSelectedTriggerId,
        filters, setFilters,
    }

    return (
        <AddTriggerDialogContext.Provider value={value}>
            {STEP_SCREENS[step]}
        </AddTriggerDialogContext.Provider>
    )
}

const STEP_SCREENS: Record<AddTriggerStep, ReactNode> = {
    mode: <ModeScreen/>,
    create: <CreateScreen/>,
    link: <LinkScreen/>,
}
