import {Button} from "@vervstack/chures"

import cls from "@/pages/tract-canvas/components/ScriptBody/ScriptBody.module.css"
import {SchemaNode, ScriptLanguage, TractStep, TractTool} from "@/processes/Tracts.ts"
import {replaceStep} from "@/processes/tractSteps.ts"
import {buildSources} from "@/pages/tract-canvas/processes/inspectorSources.ts"
import * as edits from "@/pages/tract-canvas/components/ScriptBody/processes/scriptParamEdits.ts"
import ScriptParamRow from "@/pages/tract-canvas/components/ScriptBody/components/ScriptParamRow/ScriptParamRow.tsx"
import ScriptCodeSection
    from "@/pages/tract-canvas/components/ScriptBody/components/ScriptCodeSection/ScriptCodeSection.tsx"
import Section from "@/pages/tract-canvas/components/Section/Section.tsx"
import NameField from "@/pages/tract-canvas/components/NameField/NameField.tsx"

interface Props {
    rootSteps: TractStep[]
    step: TractStep
    tools: TractTool[]
    triggerSchema?: SchemaNode
    lastOutput: unknown
    onChangeSteps: (s: TractStep[]) => void
}

export default function ScriptBody(props: Props) {
    const step = props.step
    const sources = buildSources(props.rootSteps, step.id, props.tools, props.triggerSchema)
    const inputParams = step.inputParams ?? []
    const outputParams = step.outputParams ?? []
    const lastOutput = props.lastOutput

    function update(updater: (s: TractStep) => TractStep) {
        props.onChangeSteps(replaceStep(props.rootSteps, step.id, updater))
    }

    return (
        <div className={cls.ScriptBodyContainer}>
            <Section title="Config">
                <NameField step={step} onRename={name => update(s => ({...s, name}))}/>
                <div className={cls.Row}>
                    <span className={cls.Key}>language</span>
                    <select
                        className={cls.Select}
                        value={step.language ?? ScriptLanguage.SCRIPT_LANGUAGE_JAVASCRIPT}
                        onChange={e => update(s => ({...s, language: e.target.value as ScriptLanguage}))}
                    >
                        {Object.values(ScriptLanguage)
                            .filter(l => l !== ScriptLanguage.SCRIPT_LANGUAGE_UNSPECIFIED)
                            .map(l => <option key={l} value={l}>{l.replace("SCRIPT_LANGUAGE_", "")}</option>)}
                    </select>
                </div>
            </Section>

            <Section title="Inputs">
                {inputParams.map((p, i) => (
                    <ScriptParamRow
                        key={i}
                        param={p}
                        onChange={patch => update(s => edits.updateInputParam(s, i, patch))}
                        onRemove={() => update(s => edits.removeInputParam(s, i))}
                        binding={{
                            value: step.params?.[p.name] ?? "",
                            onChange: v => update(s => ({...s, params: {...s.params, [p.name]: v}})),
                            sources,
                        }}
                    />
                ))}
                <Button variant="ghost" className={cls.AddBtn} onClick={() => update(edits.addInputParam)}>
                    + Add input param
                </Button>
            </Section>

            <Section title="Outputs">
                {outputParams.map((p, i) => (
                    <ScriptParamRow
                        key={i}
                        param={p}
                        onChange={patch => update(s => edits.updateOutputParam(s, i, patch))}
                        onRemove={() => update(s => edits.removeOutputParam(s, i))}
                    />
                ))}
                <Button variant="ghost" className={cls.AddBtn} onClick={() => update(edits.addOutputParam)}>
                    + Add output param
                </Button>
            </Section>

            <ScriptCodeSection
                name={step.name || step.id}
                inputParams={inputParams}
                outputParams={outputParams}
                code={step.code ?? ""}
                onChange={code => update(s => ({...s, code}))}
            />

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
