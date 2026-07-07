import cls from "@/pages/tract-canvas/components/ActionBody/ActionBody.module.css"

import {SchemaNode, TractStep, TractTool} from "@/processes/Tracts.ts"
import {replaceStep} from "@/processes/tractSteps.ts"
import {buildSources} from "@/pages/tract-canvas/processes/inspectorSources.ts"
import {MomCandidate} from "@/app/api/artel/mcp_keys.pb.ts"

import TemplateInput from "@/components/TemplateInput/TemplateInput.tsx"
import {connectionLabel} from "@/components/ConnectorChip/connectionLabel.ts"
import Section from "@/pages/tract-canvas/components/Section/Section.tsx"
import NameField from "@/pages/tract-canvas/components/NameField/NameField.tsx"
import {cn} from "@/app/utils/cn.ts"

const BUILTIN_MCP = "artel"

interface Props {
    rootSteps: TractStep[]
    step: TractStep
    tools: TractTool[]
    triggerSchema?: SchemaNode
    momCandidates: MomCandidate[]
    lastOutput: unknown
    onChangeSteps: (s: TractStep[]) => void
}

export default function ActionBody({rootSteps, step, tools, triggerSchema, momCandidates, lastOutput, onChangeSteps}: Props) {
    const tool = tools.find(t => t.mcp === step.mcp && t.tool === step.tool)
    const sources = buildSources(rootSteps, step.id, tools, triggerSchema)

    function update(updater: (s: TractStep) => TractStep) {
        onChangeSteps(replaceStep(rootSteps, step.id, updater))
    }

    const candidate = momCandidates.find(c => c.name === step.mcp)
    const stepConnections = candidate?.connections ?? []

    return (
        <div className={cls.ActionBodyContainer}>
            <Section title="Config">
                <NameField step={step} onRename={name => update(s => ({...s, name}))}/>
                <div className={cls.Row}>
                    <span className={cls.Key}>tool</span>
                    <span className={cls.Val}>{step.mcp}.{step.tool}</span>
                </div>
                {step.mcp !== BUILTIN_MCP && (
                    <div className={cls.Row}>
                        <span className={cls.Key}>connection</span>
                        {stepConnections.length > 0 ? (
                            <select
                                className={cn(cls.Select, !step.connection_uuid && cls.InputInvalid)}
                                value={step.connection_uuid ?? ""}
                                onChange={e => update(s => ({...s, connection_uuid: e.target.value}))}
                            >
                                <option value="">— select —</option>
                                {stepConnections.map(c => (
                                    <option key={c.id} value={c.id}>{connectionLabel(c)}</option>
                                ))}
                            </select>
                        ) : (
                            <p className={cls.Empty}>No connections available for {step.mcp}.</p>
                        )}
                    </div>
                )}
            </Section>
            {tool && Object.keys(tool.inputSchema.properties).length > 0 && (
                <Section title="Inputs">
                    {Object.entries(tool.inputSchema.properties).map(([name, def]) => (
                        <div className={cls.ParamRow} key={name}>
                            <span className={cls.Key}>{name}{tool.inputSchema.required?.includes(name) ? " *" : ""}</span>
                            {def.enum?.length ? (
                                <select
                                    className={cls.Select}
                                    value={step.params?.[name] ?? ""}
                                    onChange={e => update(s => ({...s, params: {...s.params, [name]: e.target.value}}))}
                                >
                                    <option value="">—</option>
                                    {def.enum.map(v => <option key={v} value={v}>{v}</option>)}
                                </select>
                            ) : (
                                <TemplateInput
                                    value={step.params?.[name] ?? ""}
                                    onChange={v => update(s => ({...s, params: {...s.params, [name]: v}}))}
                                    sources={sources}
                                    placeholder={def.type}
                                />
                            )}
                        </div>
                    ))}
                </Section>
            )}
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
