import {Button} from "@vervstack/chures"

import cls from "@/dialogs/ManageKeyDialog/ManageKeyDialog.module.css"
import {MomCandidate} from "@/app/api/artel/mcp_keys.pb.ts"
import DialogHead from "@/dialogs/ManageKeyDialog/components/DialogHead/DialogHead.tsx"
import ConnectionOptionList from "@/dialogs/ManageKeyDialog/components/ConnectionOptionList/ConnectionOptionList.tsx"

interface SelectConnectionScreenProps {
    candidate: MomCandidate
    selectedExternalConnectionId: string
    saving: boolean
    editing: boolean
    onSelectConnection: (id: string) => void
    onBack: () => void
    onAdd: () => void
}

export default function SelectConnectionScreen(props: SelectConnectionScreenProps) {
    const available = props.candidate.connections ?? []

    return (
        <div className={cls.Modal} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true"
             aria-labelledby="selectConnectionTitle">
            <DialogHead titleId="selectConnectionTitle" title={props.candidate.name} disabled={props.saving}/>
            <p className={cls.ModalSub}>Pick the connection this key should use for it.</p>
            <ConnectionOptionList
                available={available}
                selectedId={props.selectedExternalConnectionId}
                onSelect={props.onSelectConnection}
            />
            <div className={cls.ModalActions}>
                <Button variant="ghost" onClick={props.onBack} disabled={props.saving}>
                    Back
                </Button>
                <Button variant="primary" onClick={props.onAdd} disabled={props.saving || !props.selectedExternalConnectionId}>
                    {props.saving ? (props.editing ? "Saving…" : "Adding…") : (props.editing ? "Save" : "Add")}
                </Button>
            </div>
        </div>
    )
}
