import {Button} from "@vervstack/chures"

import cls from "@/pages/tract-canvas/components/TractCanvasInspectorBody/TractCanvasInspectorBody.module.css"
import {CanvasNode} from "@/pages/tract-canvas/processes/tractCanvasLayout.ts"
import {SchemaNode, TractStep, TractTool, TractTriggerSummary} from "@/processes/Tracts.ts"
import {Location} from "@/processes/tractSteps.ts"
import {MomCandidate} from "@/app/api/artel/mcp_keys.pb.ts"
import {cn} from "@/app/utils/cn.ts"
import {ChevronRightIcon} from "@/pages/tract-canvas/icons/ChevronRightIcon/ChevronRightIcon.tsx"
import {CloseIcon} from "@/pages/tract-canvas/icons/CloseIcon/CloseIcon.tsx"
import TriggerPanel from "@/components/TriggerPanel/TriggerPanel.tsx"
import Section from "@/pages/tract-canvas/components/Section/Section.tsx"
import OutputFields from "@/pages/tract-canvas/components/OutputFields/OutputFields.tsx"
import ActionBody from "@/pages/tract-canvas/components/ActionBody/ActionBody.tsx"
import ConditionBody from "@/pages/tract-canvas/components/ConditionBody/ConditionBody.tsx"
import ParallelBody from "@/pages/tract-canvas/components/ParallelBody/ParallelBody.tsx"
import GroupBody from "@/pages/tract-canvas/components/GroupBody/GroupBody.tsx"
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
    onToggleEnlarge: () => void
}

export default function TractCanvasInspectorBody(props: Props) {
    const {node, rootSteps, tools, triggerSchema, momCandidates, lastOutputByStepId} = props
    const {tractUuid, linkedTriggerSummaries, onChangeSteps, onOpenAddBlock, onClose} = props
    const {enlarged, onToggleEnlarge} = props
    const step = node.step

    return (
        <div className={cls.TractCanvasInspectorBodyContainer}>
            <div className={cls.Head}>
                <Button
                    variant="ghost"
                    className={cls.EnlargeBtn}
                    onClick={onToggleEnlarge}
                    aria-label={enlarged ? "Shrink inspector" : "Enlarge inspector"}
                    aria-pressed={enlarged}
                >
                    <ChevronRightIcon className={cn(cls.EnlargeChevron, enlarged && cls.EnlargeChevronOpen)}/>
                </Button>
                <div className={cls.Titles}>
                    <div className={cls.Title}>{node.kind === "trigger" ? "Trigger" : step?.name || step?.id || node.kind}</div>
                    <div className={cls.Sub}>{node.kind === "trigger" ? "trigger" : node.kind}</div>
                </div>
                <Button variant="iconDanger" className={cls.CloseBtn} onClick={onClose} aria-label="Close inspector">
                    <CloseIcon/>
                </Button>
            </div>
            <div className={cls.Body}>
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
                    <GroupBody rootSteps={rootSteps} step={step} tools={tools} triggerSchema={triggerSchema} onChangeSteps={onChangeSteps}/>
                )}
                {step && (
                    <DangerZone rootSteps={rootSteps} step={step} location={node.location} onChangeSteps={onChangeSteps} onClose={onClose}/>
                )}
                {node.kind === "action" || node.kind === "trigger" ? (
                    <Section title="Output">
                        <OutputFields schema={node.kind === "trigger" ? triggerSchema : tools.find(t => t.mcp === step?.mcp && t.tool === step?.tool)?.outputSchema}/>
                    </Section>
                ) : (
                    <Section title="Flow">
                        <Button variant="ghost" onClick={() => onOpenAddBlock(node.nextLocation, node.nextIndex)}>
                            + Add step after
                        </Button>
                    </Section>
                )}
            </div>
        </div>
    )
}
