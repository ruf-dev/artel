import {useEffect, useState} from "react"

import cls from "@/components/RunTractDialog/RunTractDialog.module.css"

import {SchemaNode, SchemaProperty, Tract, TractRun, tractsService} from "@/processes/Tracts.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"

import {Button} from "@vervstack/chures"

interface Props {
    tract?: Tract
    onRun: (params: unknown) => Promise<TractRun | undefined>
}

export default function RunTractDialog({tract, onRun}: Props) {
    const {CloseDialog} = useDialog()
    const bakeError = useBakeError()
    const [schema, setSchema] = useState<SchemaNode | null>(null)
    const [schemaLoading, setSchemaLoading] = useState(true)
    const [paramValues, setParamValues] = useState<Record<string, string>>({})
    const [rawJson, setRawJson] = useState("{}")
    const [running, setRunning] = useState(false)
    const [startedRun, setStartedRun] = useState<TractRun | null>(null)

    useEffect(() => {
        const linked = tract?.triggers ?? []
        if (linked.length === 0) {
            setSchemaLoading(false)
            return
        }

        tractsService.listTriggers()
            .then(triggers => {
                const match = triggers.find(t =>
                    linked.some(l => l.uuid === t.uuid) && Object.keys(t.payloadSchema.properties).length > 0
                )
                setSchema(match?.payloadSchema ?? null)
            })
            .catch(() => setSchema(null))
            .finally(() => setSchemaLoading(false))
    }, [tract])

    function buildParams(): unknown {
        if (schema) return coerceParams(schema.properties, paramValues)
        return JSON.parse(rawJson)
    }

    function handleRun() {
        let params: unknown
        try {
            params = buildParams()
        } catch (e) {
            bakeError("Invalid JSON", e as Error)
            return
        }
        setRunning(true)
        onRun(params)
            .then(run => setStartedRun(run ?? null))
            .catch(err => bakeError("Failed to run tract", err))
            .finally(() => setRunning(false))
    }

    if (startedRun) {
        return (
            <div className={cls.DialogContainer} role="dialog" aria-modal="true">
                <h2 className={cls.DialogTitle}>Run started</h2>
                <p className={cls.Empty}>
                    Started {formatStartedAt(startedRun.createdAt)} — follow its progress in the Runs panel.
                </p>
                <div className={cls.DialogActions}>
                    <Button variant="primary" onClick={CloseDialog}>View run in log ↓</Button>
                </div>
            </div>
        )
    }

    return (
        <div className={cls.DialogContainer} role="dialog" aria-modal="true">
            <h2 className={cls.DialogTitle}>Run tract</h2>
            {schemaLoading ? (
                <p className={cls.Empty}>Loading…</p>
            ) : schema ? (
                <ParamsList schema={schema} values={paramValues} onChange={(k, v) => setParamValues(p => ({...p, [k]: v}))}/>
            ) : (
                <textarea
                    className={cls.JsonInput}
                    value={rawJson}
                    onChange={e => setRawJson(e.target.value)}
                    rows={8}
                    spellCheck={false}
                />
            )}
            <div className={cls.DialogActions}>
                <Button variant="ghost" onClick={CloseDialog} disabled={running}>Cancel</Button>
                <Button variant="primary" onClick={handleRun} disabled={running}>
                    {running ? "Running…" : "Run"}
                </Button>
            </div>
        </div>
    )
}

function formatStartedAt(iso: string | undefined): string {
    const d = iso ? new Date(iso) : null
    if (!d || isNaN(d.getTime())) return "just now"
    return d.toLocaleTimeString(undefined, {hour: "2-digit", minute: "2-digit", second: "2-digit"})
}

function ParamsList({schema, values, onChange}: {
    schema: SchemaNode
    values: Record<string, string>
    onChange: (name: string, value: string) => void
}) {
    const entries = Object.entries(schema.properties)
    if (entries.length === 0) return null
    return (
        <div className={cls.ParamsSection}>
            {entries.map(([name, def]) => (
                <ParamRow
                    key={name}
                    name={name}
                    def={def}
                    required={schema.required?.includes(name) ?? false}
                    value={values[name] ?? ""}
                    onChange={onChange}
                />
            ))}
        </div>
    )
}

function ParamRow({name, def, required, value, onChange}: {
    name: string
    def: SchemaProperty
    required: boolean
    value: string
    onChange: (name: string, value: string) => void
}) {
    return (
        <div className={cls.ParamRow}>
            <span className={cls.ParamName}>{name}{required ? " *" : ""}</span>
            {def.description && <span className={cls.ParamDesc}>{def.description}</span>}
            <ParamInput def={def} value={value} onChange={v => onChange(name, v)}/>
        </div>
    )
}

function ParamInput({def, value, onChange}: {
    def: SchemaProperty
    value: string
    onChange: (value: string) => void
}) {
    if (def.enum?.length) {
        return (
            <select className={cls.ParamField} value={value} onChange={e => onChange(e.target.value)}>
                <option value="">—</option>
                {def.enum.map(v => <option key={v} value={v}>{v}</option>)}
            </select>
        )
    }
    if (def.type === "integer" || def.type === "number") {
        return <input className={cls.ParamField} type="number" value={value} onChange={e => onChange(e.target.value)}/>
    }
    if (def.type === "boolean") {
        return <input type="checkbox" checked={value === "true"} onChange={e => onChange(String(e.target.checked))}/>
    }
    return <input className={cls.ParamField} type="text" value={value} onChange={e => onChange(e.target.value)}/>
}

function coerceParams(props: Record<string, SchemaProperty>, values: Record<string, string>): Record<string, unknown> {
    const out: Record<string, unknown> = {}
    for (const [name, def] of Object.entries(props)) {
        const raw = values[name]
        if (raw === undefined || raw === "") continue
        if (def.type === "integer" || def.type === "number") {
            const n = Number(raw)
            out[name] = Number.isNaN(n) ? raw : n
        } else if (def.type === "boolean") {
            out[name] = raw === "true"
        } else if (def.type === "array" || def.type === "object") {
            try {
                out[name] = JSON.parse(raw)
            } catch {
                out[name] = raw
            }
        } else {
            out[name] = raw
        }
    }
    return out
}
