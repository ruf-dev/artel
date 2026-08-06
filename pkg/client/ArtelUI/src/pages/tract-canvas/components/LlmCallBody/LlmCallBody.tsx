import {Input} from "@vervstack/chures"

import cls from "@/pages/tract-canvas/components/LlmCallBody/LlmCallBody.module.css"
import {SchemaNode, TractStep} from "@/processes/Tracts.ts"
import {replaceStep} from "@/processes/tractSteps.ts"
import {buildSources} from "@/pages/tract-canvas/processes/inspectorSources.ts"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import TemplateInput from "@/components/TemplateInput/TemplateInput.tsx"
import {connectionLabel} from "@/components/ConnectorChip/connectionLabel.ts"
import Section from "@/pages/tract-canvas/components/Section/Section.tsx"
import NameField from "@/pages/tract-canvas/components/NameField/NameField.tsx"
import {cn} from "@/app/utils/cn.ts"

interface Props {
    rootSteps: TractStep[]
    step: TractStep
    triggerSchema?: SchemaNode
    lastOutput: unknown
    onChangeSteps: (s: TractStep[]) => void
}

export default function LlmCallBody(props: Props) {
    const step = props.step
    const sources = buildSources(props.rootSteps, step.id, [], props.triggerSchema)
    const lastOutput = props.lastOutput
    const {connections} = useExternalConnections()

    function update(updater: (s: TractStep) => TractStep) {
        props.onChangeSteps(replaceStep(props.rootSteps, step.id, updater))
    }

    const LLM_PROVIDERS = [ExternalProvider.EXTERNAL_PROVIDER_ANTHROPIC, ExternalProvider.EXTERNAL_PROVIDER_OPENAI]
    const llmConnections = connections.filter(c => c.provider && LLM_PROVIDERS.includes(c.provider))
    const selectedConnection = llmConnections.find(c => c.id === step.llmConnectionId)
    const availableModels = selectedConnection?.generic?.fields?.available_models
        ? selectedConnection.generic.fields.available_models.split(",").filter(Boolean)
        : []

    return (
        <div className={cls.LlmCallBodyContainer}>
            <Section title="Config">
                <NameField step={step} onRename={name => update(s => ({...s, name}))}/>
                <div className={cls.Row}>
                    <span className={cls.Key}>connection</span>
                    {llmConnections.length > 0 ? (
                        <select
                            className={cn(cls.Select, !step.llmConnectionId && cls.InputInvalid)}
                            value={step.llmConnectionId ?? ""}
                            onChange={e => update(s => ({...s, llmConnectionId: e.target.value, llmModel: undefined}))}
                        >
                            <option value="">— select —</option>
                            {llmConnections.map(c => (
                                <option key={c.id} value={c.id}>{connectionLabel(c)}</option>
                            ))}
                        </select>
                    ) : (
                        <p className={cls.Empty}>No LLM connections available.</p>
                    )}
                </div>
                <div className={cls.Row}>
                    <span className={cls.Key}>model</span>
                    {availableModels.length > 0 ? (
                        <select
                            className={cn(cls.Select, !step.llmModel && cls.InputInvalid)}
                            value={step.llmModel ?? ""}
                            onChange={e => update(s => ({...s, llmModel: e.target.value}))}
                        >
                            <option value="">— select —</option>
                            {availableModels.map(m => <option key={m} value={m}>{m}</option>)}
                        </select>
                    ) : (
                        <p className={cls.Empty}>No models available — check the connection.</p>
                    )}
                </div>
                <div className={cls.Row}>
                    <span className={cls.Key}>max tokens</span>
                    <Input
                        className={cls.NumberInputWrapper}
                        inputClassName={cls.NumberInput}
                        type="number"
                        value={step.maxTokens !== undefined ? String(step.maxTokens) : ""}
                        setValue={v => update(s => ({...s, maxTokens: v === "" ? undefined : Number(v)}))}
                    />
                </div>
            </Section>
            <Section title="Prompt">
                <div className={cls.ParamRow}>
                    <span className={cls.Key}>system prompt</span>
                    <TemplateInput
                        value={step.systemPrompt ?? ""}
                        onChange={v => update(s => ({...s, systemPrompt: v}))}
                        sources={sources}
                        placeholder="optional"
                    />
                </div>
                <div className={cls.ParamRow}>
                    <span className={cls.Key}>prompt</span>
                    <TemplateInput
                        value={step.prompt ?? ""}
                        onChange={v => update(s => ({...s, prompt: v}))}
                        sources={sources}
                        placeholder="prompt"
                    />
                </div>
            </Section>
            {lastOutput !== undefined && lastOutput !== null && typeof lastOutput === "object" && (
                <Section title="Last output">
                    {Object.entries(lastOutput as Record<string, unknown>).map(([k, v]) => (
                        <div className={cls.Row} key={k}>
                            <span className={cls.Key}>{k}</span>
                            <span className={cls.ValDim}>{String(v)}</span>
                        </div>
                    ))}
                </Section>
            )}
        </div>
    )
}
