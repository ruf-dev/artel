import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import {connectionLabel} from "@/components/ConnectorChip/connectionLabel.ts"
import CardHeader from "@/components/TractStepTree/components/CardHeader.tsx"
import cls from "@/components/TractStepTree/TractStepTree.module.css"
import {TractStep} from "@/processes/Tracts.ts"

const PROMPT_PREVIEW_LENGTH = 80

interface Props {
    step: TractStep
    onUpdate: (updater: (s: TractStep) => TractStep) => void
    onDelete: () => void
}

/** Read-only tree-view card for an "llm_call" step — a leaf step (no nested branches/lanes), so
 * it needs its own card rather than falling through TractStepTree's default "group" branch,
 * which would render an empty add-child list that doesn't apply here. */
export default function LlmCallCard({step, onUpdate, onDelete}: Props) {
    const {connections} = useExternalConnections()
    const conn = connections.find(
        c => c.id === step.llmConnectionId && c.provider === ExternalProvider.EXTERNAL_PROVIDER_ANTHROPIC,
    )
    const prompt = step.prompt ?? ""
    const truncated = prompt.length > PROMPT_PREVIEW_LENGTH
    const promptPreview = truncated ? `${prompt.slice(0, PROMPT_PREVIEW_LENGTH)}…` : prompt

    return (
        <div className={cls.Card}>
            <CardHeader
                step={step}
                onUpdate={onUpdate}
                onDelete={onDelete}
                right={
                    <>
                        {step.llmModel && <span className={cls.ConnChip}>{step.llmModel}</span>}
                        {conn && <span className={cls.ConnChip}>{connectionLabel(conn)}</span>}
                    </>
                }
            />
            {prompt && (
                <div className={cls.ParamsSection}>
                    <div className={cls.ParamRow}>
                        <span className={cls.ParamName}>prompt</span>
                        <span className={cls.ParamDesc}>{promptPreview}</span>
                    </div>
                </div>
            )}
        </div>
    )
}
