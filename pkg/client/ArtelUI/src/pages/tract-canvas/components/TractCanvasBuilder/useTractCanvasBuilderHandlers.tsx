import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import {useTracts} from "@/app/hooks/Tracts.ts"
import TractBlockPicker from "@/pages/tract-canvas/components/TractBlockPicker/TractBlockPicker.tsx"
import InsertConflictDialog from "@/pages/tract-canvas/dialogs/InsertConflictDialog/InsertConflictDialog.tsx"
import {StepDraft} from "@/components/StepPickerDialog/StepPickerDialog.tsx"
import RunTractDialog from "@/components/RunTractDialog/RunTractDialog.tsx"
import {buildStepFromDraft, collectAllStepIds, insertBlockAfter, Location} from "@/processes/tractSteps.ts"
import {DropSide, findStepById, moveStepRelativeTo, wouldConflict} from "@/processes/tractStepsMove.ts"
import {Tract, TractDefinition, TractRun} from "@/processes/Tracts.ts"

interface Deps {
    tract: Tract
    name: string
    definition: TractDefinition
    setDefinition: (updater: TractDefinition | ((d: TractDefinition) => TractDefinition)) => void
    setSaving: (saving: boolean) => void
    setWarnings: (warnings: string[]) => void
    setSavedName: (name: string) => void
    setSavedDefinition: (definition: TractDefinition) => void
    startRun: (params: unknown) => Promise<TractRun | undefined>
}

export function useTractCanvasBuilderHandlers(deps: Deps) {
    const {updateTract} = useTracts()
    const bakeError = useBakeError()
    const {OpenDialog} = useDialog()

    function openRunDialog() {
        OpenDialog(<RunTractDialog tract={deps.tract} onRun={deps.startRun}/>)
    }

    function handleSave() {
        deps.setSaving(true)
        updateTract(deps.tract.uuid, deps.name, deps.tract.description, deps.definition)
            .then(({warnings}) => {
                deps.setWarnings(warnings)
                deps.setSavedName(deps.name)
                deps.setSavedDefinition(deps.definition)
            })
            .catch(err => bakeError("Failed to save tract", err))
            .finally(() => deps.setSaving(false))
    }

    function handleChangeSteps(steps: TractDefinition["steps"]) {
        deps.setDefinition({steps})
    }

    function openAddBlock(location: Location, index: number) {
        OpenDialog(
            <TractBlockPicker
                onConfirm={(draft: StepDraft) => {
                    const existingIds = collectAllStepIds(deps.definition.steps)
                    const newStep = buildStepFromDraft(draft, existingIds)
                    existingIds.add(newStep.id)
                    deps.setDefinition(d => ({steps: insertBlockAfter(d.steps, location, index, newStep, existingIds)}))
                }}
            />
        )
    }

    function handleMoveStep(sourceId: string, targetId: string, side: DropSide) {
        if (side === "before" || !wouldConflict(deps.definition.steps, sourceId, targetId, side)) {
            deps.setDefinition(d => ({steps: moveStepRelativeTo(d.steps, sourceId, targetId, side)}))
            return
        }

        const sourceName = findStepById(deps.definition.steps, sourceId)?.name || sourceId
        const targetName = findStepById(deps.definition.steps, targetId)?.name || "Trigger"
        OpenDialog(
            <InsertConflictDialog
                sourceName={sourceName}
                targetName={targetName}
                onChoose={choice => deps.setDefinition(d => (
                    {steps: moveStepRelativeTo(d.steps, sourceId, targetId, side, choice)}
                ))}
            />
        )
    }

    return {openRunDialog, handleSave, handleChangeSteps, openAddBlock, handleMoveStep}
}
