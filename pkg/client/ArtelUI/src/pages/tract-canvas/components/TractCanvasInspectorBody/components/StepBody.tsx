import {CanvasNode} from "@/pages/tract-canvas/processes/tractCanvasLayout.ts"
import {SchemaNode, TractStep, TractTool, TractTriggerSummary} from "@/processes/Tracts.ts"
import {Location} from "@/processes/tractSteps.ts"
import {MomCandidate} from "@/app/api/artel/mcp_keys.pb.ts"
import TriggerPanel from "@/components/TriggerPanel/TriggerPanel.tsx"
import ActionBody from "@/pages/tract-canvas/components/ActionBody/ActionBody.tsx"
import ConditionBody from "@/pages/tract-canvas/components/ConditionBody/ConditionBody.tsx"
import ParallelBody from "@/pages/tract-canvas/components/ParallelBody/ParallelBody.tsx"
import GroupBody from "@/pages/tract-canvas/components/GroupBody/GroupBody.tsx"
import ScriptBody from "@/pages/tract-canvas/components/ScriptBody/ScriptBody.tsx"
import DangerZone from "@/pages/tract-canvas/components/DangerZone/DangerZone.tsx"

interface Props {
    node: CanvasNode
    rootSteps: TractStep[]
    tools: TractTool[]
    triggerSchema?: SchemaNode
    momCandidates: MomCandidate[]
    lastOutputByStepId: Record<string, unknown>
    tractUuid: string
    linkedTriggerSummaries: TractTriggerSummary[]
    onChangeSteps: (newRootSteps: TractStep[]) => void
    onOpenAddBlock: (location: Location, index: number) => void
    onClose: () => void
    enlarged: boolean
}

/** StepBody renders whichever kind-specific body (+ the shared DangerZone) matches
 * props.node.kind — split out of TractCanvasInspectorBody purely to keep that file under
 * the project's max-lines-per-function limit. */
export default function StepBody(props: Props) {
    const {node, rootSteps, tools, triggerSchema, momCandidates, lastOutputByStepId} = props
    const {tractUuid, linkedTriggerSummaries, onChangeSteps, onOpenAddBlock, onClose, enlarged} = props
    const step = node.step

    return (
        <>
            {node.kind === "trigger" && (
                <TriggerPanel tractUuid={tractUuid} linkedTriggerSummaries={linkedTriggerSummaries}/>
            )}
            {step && node.kind === "action" && (
                <ActionBody
                    rootSteps={rootSteps} step={step} tools={tools} triggerSchema={triggerSchema}
                    momCandidates={momCandidates}
                    lastOutput={lastOutputByStepId[step.id]}
                    onChangeSteps={onChangeSteps}
                />
            )}
            {step && node.kind === "condition" && (
                <ConditionBody
                    rootSteps={rootSteps} step={step} tools={tools} triggerSchema={triggerSchema}
                    onChangeSteps={onChangeSteps} onOpenAddBlock={onOpenAddBlock} enlarged={enlarged}
                />
            )}
            {step && node.kind === "parallel" && (
                <ParallelBody step={step} onOpenAddBlock={onOpenAddBlock}/>
            )}
            {step && node.kind === "group" && (
                <GroupBody
                    rootSteps={rootSteps}
                    step={step}
                    tools={tools}
                    triggerSchema={triggerSchema}
                    onChangeSteps={onChangeSteps}
                />
            )}
            {step && node.kind === "script" && (
                <ScriptBody
                    rootSteps={rootSteps} step={step} tools={tools} triggerSchema={triggerSchema}
                    lastOutput={lastOutputByStepId[step.id]}
                    onChangeSteps={onChangeSteps}
                />
            )}
            {step && (
                <DangerZone
                    rootSteps={rootSteps}
                    step={step}
                    location={node.location}
                    onChangeSteps={onChangeSteps}
                    onClose={onClose}
                />
            )}
        </>
    )
}
