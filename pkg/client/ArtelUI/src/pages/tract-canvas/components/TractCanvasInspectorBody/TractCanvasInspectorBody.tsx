import {Button} from "@vervstack/chures"

import cls from "@/pages/tract-canvas/components/TractCanvasInspectorBody/TractCanvasInspectorBody.module.css"
import {CanvasNode} from "@/pages/tract-canvas/processes/tractCanvasLayout.ts"
import {SchemaNode, TractStep, TractTool, TractTriggerSummary} from "@/processes/Tracts.ts"
import {Location} from "@/processes/tractSteps.ts"
import {MomCandidate} from "@/app/api/artel/mcp_keys.pb.ts"
import {cn} from "@/app/utils/cn.ts"
import {ChevronRightIcon} from "@/pages/tract-canvas/icons/ChevronRightIcon/ChevronRightIcon.tsx"
import {CloseIcon} from "@/pages/tract-canvas/icons/CloseIcon/CloseIcon.tsx"
import JsonBlock from "@/components/JsonBlock/JsonBlock.tsx"
import Section from "@/pages/tract-canvas/components/Section/Section.tsx"
import OutputFields from "@/pages/tract-canvas/components/OutputFields/OutputFields.tsx"
import StepBody from "@/pages/tract-canvas/components/TractCanvasInspectorBody/components/StepBody.tsx"

// Read-only value actually used the last time this node ran, shown only while a run is
// selected in the log panel — the editable inputs elsewhere in the sidebar are unaffected.
function runValueFor(node: CanvasNode, activeTriggerPayload: unknown, activeInputByStepId: Record<string, unknown>) {
    if (node.kind === "trigger") return activeTriggerPayload
    if ((node.kind === "action" || node.kind === "script") && node.step) return activeInputByStepId[node.step.id]
    return undefined
}

interface Props {
    node: CanvasNode
    rootSteps: TractStep[]
    tools: TractTool[]
    triggerSchema?: SchemaNode
    momCandidates: MomCandidate[]
    lastOutputByStepId: Record<string, unknown>
    activeInputByStepId: Record<string, unknown>
    activeTriggerPayload: unknown
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
    const {activeInputByStepId, activeTriggerPayload} = props
    const {tractUuid, linkedTriggerSummaries, onChangeSteps, onOpenAddBlock, onClose} = props
    const {enlarged, onToggleEnlarge} = props
    const step = node.step
    const runValue = runValueFor(node, activeTriggerPayload, activeInputByStepId)
    const runValueLabel = node.kind === "trigger" ? "Trigger payload" : "Step input"

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
                    <div className={cls.Title}>
                        {node.kind === "trigger" ? "Trigger" : step?.name || step?.id || node.kind}
                    </div>
                    <div className={cls.Sub}>{node.kind === "trigger" ? "trigger" : node.kind}</div>
                </div>
                <Button variant="iconDanger" className={cls.CloseBtn} onClick={onClose} aria-label="Close inspector">
                    <CloseIcon/>
                </Button>
            </div>
            <div className={cls.Body}>
                <StepBody
                    node={node} rootSteps={rootSteps} tools={tools} triggerSchema={triggerSchema}
                    momCandidates={momCandidates} lastOutputByStepId={lastOutputByStepId}
                    tractUuid={tractUuid} linkedTriggerSummaries={linkedTriggerSummaries}
                    onChangeSteps={onChangeSteps} onOpenAddBlock={onOpenAddBlock} onClose={onClose}
                    enlarged={enlarged}
                />
                {node.kind === "action" || node.kind === "trigger" ? (
                    <Section title="Output">
                        <OutputFields
                            schema={
                                node.kind === "trigger"
                                    ? triggerSchema
                                    : tools.find(t => t.mcp === step?.mcp && t.tool === step?.tool)?.outputSchema
                            }
                        />
                    </Section>
                ) : node.kind === "script" ? null : (
                    <Section title="Flow">
                        <Button
                            variant="ghost"
                            className={cls.AddStepBtn}
                            onClick={() => onOpenAddBlock(node.nextLocation, node.nextIndex)}
                        >
                            + Add step after
                        </Button>
                    </Section>
                )}
                {runValue !== undefined && (
                    <Section title="Run values">
                        <JsonBlock label={runValueLabel} value={runValue}/>
                    </Section>
                )}
            </div>
        </div>
    )
}
