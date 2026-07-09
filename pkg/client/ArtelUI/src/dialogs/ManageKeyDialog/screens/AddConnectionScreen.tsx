import {Button} from "@vervstack/chures"

import cls from "@/dialogs/ManageKeyDialog/ManageKeyDialog.module.css"
import {McpConnectorInfo, MomCandidate} from "@/app/api/artel/mcp_keys.pb.ts"
import DialogHead from "@/dialogs/ManageKeyDialog/components/DialogHead/DialogHead.tsx"
import CandidateOptionList from "@/dialogs/ManageKeyDialog/components/CandidateOptionList/CandidateOptionList.tsx"

interface AddConnectionScreenProps {
    connectors: McpConnectorInfo[]
    onSelect: (candidate: MomCandidate) => void
    onBack: () => void
}

export default function AddConnectionScreen({connectors, onSelect, onBack}: AddConnectionScreenProps) {
    return (
        <div className={cls.Modal} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true"
             aria-labelledby="addConnectionTitle">
            <DialogHead titleId="addConnectionTitle" title="Add connection"/>
            <p className={cls.ModalSub}>Pick a service to connect to this key.</p>
            <CandidateOptionList connectors={connectors} onSelect={onSelect}/>
            <div className={cls.ModalActions}>
                <Button variant="ghost" onClick={onBack}>Back</Button>
            </div>
        </div>
    )
}
