import {useEffect, useState} from "react"

import cls from "@/pages/tract-canvas/widgets/TractCanvasInspector/TractCanvasInspector.module.css"

import {CanvasNode} from "@/pages/tract-canvas/processes/tractCanvasLayout.ts"
import {SchemaNode, TractStep, TractTool, TractTriggerSummary} from "@/processes/Tracts.ts"
import {Location} from "@/processes/tractSteps.ts"
import {MomCandidate} from "@/app/api/artel/mcp_keys.pb.ts"
import {cn} from "@/app/utils/cn.ts"

import TractCanvasInspectorBody from "@/pages/tract-canvas/components/TractCanvasInspectorBody/TractCanvasInspectorBody.tsx"

interface Props {
    node: CanvasNode | null
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
}

export default function TractCanvasInspector({node, rootSteps, tools, triggerSchema, momCandidates, lastOutputByStepId, tractUuid, linkedTriggerSummaries, onChangeSteps, onOpenAddBlock, onClose}: Props) {
    const [enlarged, setEnlarged] = useState(false)

    useEffect(() => {
        setEnlarged(false)
    }, [node?.id])

    return (
        <div className={cls.TractCanvasInspectorContainer}>
            {node && enlarged && (
                <div className={cls.Backdrop} onClick={() => setEnlarged(false)}/>
            )}
            <div className={cn(cls.Panel, node && cls.Open, node && enlarged && cls.Enlarged)}>
                {node && (
                    <TractCanvasInspectorBody
                        key={node.id}
                        node={node}
                        rootSteps={rootSteps}
                        tools={tools}
                        triggerSchema={triggerSchema}
                        momCandidates={momCandidates}
                        lastOutputByStepId={lastOutputByStepId}
                        tractUuid={tractUuid}
                        linkedTriggerSummaries={linkedTriggerSummaries}
                        onChangeSteps={onChangeSteps}
                        onOpenAddBlock={onOpenAddBlock}
                        onClose={onClose}
                        enlarged={enlarged}
                        onToggleEnlarge={() => setEnlarged(e => !e)}
                    />
                )}
            </div>
        </div>
    )
}
