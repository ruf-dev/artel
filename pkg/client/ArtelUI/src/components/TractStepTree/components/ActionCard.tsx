import {useState} from "react"

import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {connectionLabel} from "@/components/ConnectorChip/connectionLabel.ts"
import CardHeader from "@/components/TractStepTree/components/CardHeader.tsx"
import ParamRow from "@/components/TractStepTree/components/ParamRow.tsx"
import SchemaTree from "@/components/TractStepTree/components/SchemaTree.tsx"
import {buildSourcesFor} from "@/components/TractStepTree/processes/tractStepTreeHelpers.ts"
import cls from "@/components/TractStepTree/TractStepTree.module.css"
import {SchemaNode, TractStep, TractTool} from "@/processes/Tracts.ts"

interface Props {
    rootSteps: TractStep[]
    step: TractStep
    tools: TractTool[]
    triggerSchema?: SchemaNode
    onUpdate: (updater: (s: TractStep) => TractStep) => void
    onDelete: () => void
}

export default function ActionCard({rootSteps, step, tools, triggerSchema, onUpdate, onDelete}: Props) {
    const {connections} = useExternalConnections()
    const [outputOpen, setOutputOpen] = useState(false)

    const tool = tools.find(t => t.mcp === step.mcp && t.tool === step.tool)
    const conn = connections.find(c => c.id === step.connection_uuid)
    const sources = buildSourcesFor(rootSteps, step.id, tools, triggerSchema)

    function setParam(name: string, value: string) {
        onUpdate(s => ({...s, params: {...s.params, [name]: value}}))
    }

    return (
        <div className={cls.Card}>
            <CardHeader
                step={step}
                onUpdate={onUpdate}
                onDelete={onDelete}
                right={
                    <>
                        <span className={cls.ConnChip}>{step.mcp}.{step.tool}</span>
                        {conn && <span className={cls.ConnChip}>{connectionLabel(conn)}</span>}
                    </>
                }
            />
            {tool && Object.keys(tool.inputSchema.properties).length > 0 && (
                <div className={cls.ParamsSection}>
                    {Object.entries(tool.inputSchema.properties).map(([name, def]) => (
                        <ParamRow
                            key={name}
                            name={name}
                            def={def}
                            required={tool.inputSchema.required?.includes(name) ?? false}
                            value={step.params?.[name] ?? ""}
                            sources={sources}
                            onChange={setParam}
                        />
                    ))}
                </div>
            )}
            {tool && Object.keys(tool.outputSchema.properties).length > 0 && (
                <div className={cls.Collapsible}>
                    <div className={cls.CollapsibleHeader} onClick={() => setOutputOpen(o => !o)}>
                        <span>{outputOpen ? "▾" : "▸"} Output</span>
                    </div>
                    {outputOpen && <SchemaTree schema={tool.outputSchema}/>}
                </div>
            )}
        </div>
    )
}
