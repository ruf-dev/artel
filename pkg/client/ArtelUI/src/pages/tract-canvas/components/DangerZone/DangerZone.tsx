import {Button, ConfirmDialog} from "@vervstack/chures"

import cls from "@/pages/tract-canvas/components/DangerZone/DangerZone.module.css"
import {TractStep} from "@/processes/Tracts.ts"
import {hasChildren, Location, removeStepAt} from "@/processes/tractSteps.ts"
import {useDialog} from "@/app/hooks/Dialog"
import Section from "@/pages/tract-canvas/components/Section/Section.tsx"

interface Props {
    rootSteps: TractStep[]
    step: TractStep
    location: Location
    onChangeSteps: (s: TractStep[]) => void
    onClose: () => void
}

export default function DangerZone({rootSteps, step, location, onChangeSteps, onClose}: Props) {
    const {OpenDialog, CloseDialog} = useDialog()

    function doDelete() {
        onChangeSteps(removeStepAt(rootSteps, location, step.id))
        onClose()
    }

    function handleDelete() {
        if (hasChildren(step)) {
            OpenDialog(
                <ConfirmDialog
                    title="Delete step"
                    message={`Delete "${step.name || step.id}"? Its nested branches/lanes will be removed too.`}
                    confirmLabel="Delete"
                    danger
                    onClose={CloseDialog}
                    onConfirm={async () => doDelete()}
                />
            )
            return
        }
        doDelete()
    }

    return (
        <div className={cls.DangerZoneContainer}>
            <Section title="Danger zone">
                <Button variant="danger" onClick={handleDelete}>Delete step</Button>
            </Section>
        </div>
    )
}
