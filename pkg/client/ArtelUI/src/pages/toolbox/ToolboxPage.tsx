import {useEffect, useState} from "react"
import {useNavigate} from "react-router-dom"

import cls from "@/pages/toolbox/ToolboxPage.module.css"

import {ImapOperation, ImapToolAction, McpToolInfo, MomCandidate, SmtpOperation, SmtpToolAction, ToolParamDef} from "@/app/api/artel/mcp_keys.pb.ts"
import {ExternalConnectionInfo} from "@/app/api/artel/external_connections.pb.ts"
import {Path} from "@/app/routing/Router.tsx"
import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import {useBakeError} from "@/app/hooks/useErrorToast"
import useUser from "@/hooks/user/User.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {Button} from "@vervstack/chures"
import SelectOption from "@/components/SelectOption/SelectOption.tsx"
import {connectionLabel} from "@/components/ConnectorChip/connectionLabel.ts"

export default function ToolboxPage() {
    const navigate = useNavigate()
    const {auth} = useUser()
    const {fetchMomCandidates} = useMcpKeys()

    useEffect(() => {
        if (!auth.isAuthenticated()) {
            navigate(Path.InitPage)
        }
    }, [auth, navigate])

    useEffect(() => {
        if (auth.isAuthenticated()) {
            void fetchMomCandidates()
        }
    }, [auth, fetchMomCandidates])

    return (
        <div className={cls.Root}>
            <HeroSegment/>
            <ContentSegment/>
        </div>
    )
}

function HeroSegment() {
    const {momCandidates, momCandidatesLoading} = useMcpKeys()

    return (
        <div className={cls.Hero}>
            <div className={cls.HeroTitles}>
                <div className={cls.Eyebrow}>MCP</div>
                <h1 className={cls.HeroTitle}>Toolbox</h1>
                <p className={cls.HeroSub}>
                    <b>{momCandidatesLoading ? "…" : `${momCandidates.length} ${momCandidates.length === 1 ? "tool" : "tools"}`}</b>
                    {" · "}<span>available in this installation</span>
                </p>
            </div>
        </div>
    )
}

function ContentSegment() {
    const {momCandidates, momCandidatesLoading} = useMcpKeys()

    if (momCandidatesLoading) {
        return (
            <div className={cls.Content}>
                <p className={cls.Empty}>Loading…</p>
            </div>
        )
    }

    return (
        <div className={cls.Content}>
            <div className={cls.List}>
                {momCandidates.map(c => (
                    <MomCandidateRow key={c.name} candidate={c}/>
                ))}
                {momCandidates.length === 0 && (
                    <p className={cls.Empty}>No tools available.</p>
                )}
            </div>
        </div>
    )
}

function MomCandidateRow({candidate}: { candidate: MomCandidate }) {
    const {OpenDialog, SetClosable} = useDialog()
    const connCount = candidate.connections?.length ?? 0
    const toolCount = candidate.tools?.length ?? 0

    function onToolboxClick() {
        OpenDialog(<ToolsDialog candidate={candidate}/>)
        SetClosable(true)
    }

    return (
        <div className={cls.Row} onClick={onToolboxClick} role="button" tabIndex={0}>
            <div className={cls.RowHeader}>
                <span className={cls.RowName}>{candidate.name}</span>
                {toolCount > 0 && (
                    <span className={cls.Chip}>{toolCount} {toolCount === 1 ? "tool" : "tools"}</span>
                )}
                {connCount > 0 && (
                    <span className={cls.Chip}>{connCount} {connCount === 1 ? "connection" : "connections"}</span>
                )}
            </div>
            <div className={cls.RowAuthor}>{candidate.author}</div>
            <div className={cls.RowDesc}>{candidate.description}</div>
        </div>
    )
}

function ToolsDialog({candidate}: { candidate: MomCandidate }) {
    const tools = candidate.tools ?? []
    const [activeTool, setActiveTool] = useState<McpToolInfo | null>(null)

    if (activeTool) {
        return <ToolDetail candidate={candidate} tool={activeTool} onBack={() => setActiveTool(null)}/>
    }

    return (
        <div className={cls.DialogContainer} role="dialog" aria-modal="true">
            <h2 className={cls.DialogTitle}>{candidate.name}</h2>
            <p className={cls.DialogDesc}>{candidate.description}</p>
            <div className={cls.ToolList}>
                {tools.map(t => (
                    <ToolRow key={t.name} tool={t} onClick={() => setActiveTool(t)}/>
                ))}
                {tools.length === 0 && <p className={cls.ToolEmpty}>No tools defined.</p>}
            </div>
        </div>
    )
}

function ToolDetail({candidate, tool, onBack}: { candidate: MomCandidate; tool: McpToolInfo; onBack: () => void }) {
    const {executeMomTool} = useMcpKeys()
    const bakeError = useBakeError()
    const connections = candidate.connections ?? []
    const [selectedConnectionId, setSelectedConnectionId] = useState<string>(connections.length === 1 ? (connections[0].id ?? "") : "")
    const [paramValues, setParamValues] = useState<Record<string, string>>({})
    const [running, setRunning] = useState(false)
    const [result, setResult] = useState<string | null>(null)

    function setParam(name: string, value: string) {
        setParamValues(prev => ({...prev, [name]: value}))
    }

    function onRun() {
        if (!selectedConnectionId) return
        setRunning(true)
        setResult(null)
        const params = coerceParams(tool.params ?? {}, paramValues)
        executeMomTool(candidate.name ?? "", tool.name ?? "", selectedConnectionId, params)
            .then(setResult)
            .catch(err => bakeError("Failed to run tool", err))
            .finally(() => setRunning(false))
    }

    return (
        <div className={cls.DialogContainer} role="dialog" aria-modal="true">
            <h2 className={cls.DialogTitle}>{tool.name}</h2>
            <p className={cls.DialogDesc}>{tool.description}</p>
            <div className={cls.ToolDetail}>
                {tool.smtp && <SmtpActionView action={tool.smtp}/>}
                {tool.imap && <ImapActionView action={tool.imap}/>}
                <ConnectionPicker connections={connections} selectedId={selectedConnectionId} onSelect={setSelectedConnectionId}/>
                {tool.params && <ParamsList params={tool.params} values={paramValues} onChange={setParam}/>}
                {result !== null && <ResultPanel result={result}/>}
            </div>
            <div className={cls.DialogActions}>
                <Button variant="ghost" onClick={onBack}>Back</Button>
                <Button variant="primary" disabled={running || !selectedConnectionId} onClick={onRun}>
                    {running ? "Running…" : "Run"}
                </Button>
            </div>
        </div>
    )
}

function coerceParams(defs: Record<string, ToolParamDef>, values: Record<string, string>): Record<string, unknown> {
    const out: Record<string, unknown> = {}
    for (const [name, def] of Object.entries(defs)) {
        const raw = values[name]
        if (raw === undefined || raw === "") continue
        if (def.integerParam) {
            const n = Number(raw)
            out[name] = Number.isNaN(n) ? raw : n
        } else {
            out[name] = raw
        }
    }
    return out
}

function ToolRow({tool, onClick}: { tool: McpToolInfo; onClick: () => void }) {
    return (
        <div className={cls.ToolRow} onClick={onClick} role="button" tabIndex={0}>
            <span className={cls.ToolName}>{tool.name}</span>
            <p className={cls.ToolDesc}>{tool.description}</p>
        </div>
    )
}

function SmtpActionView({action}: { action: SmtpToolAction }) {
    return (
        <div className={cls.ActionContainer}>
            <label className={cls.ActionLabel}>Operation</label>
            <select className={cls.ActionSelect} disabled value={action.operation ?? SmtpOperation.SMTP_OP_UNSPECIFIED}>
                <option value={SmtpOperation.SMTP_OP_SEND}>send</option>
            </select>
        </div>
    )
}

function ImapActionView({action}: { action: ImapToolAction }) {
    return (
        <div className={cls.ActionContainer}>
            <label className={cls.ActionLabel}>Operation</label>
            <select className={cls.ActionSelect} disabled value={action.operation ?? ImapOperation.IMAP_OP_UNSPECIFIED}>
                <option value={ImapOperation.IMAP_OP_LIST_FOLDERS}>list folders</option>
                <option value={ImapOperation.IMAP_OP_LIST_MESSAGES}>list messages</option>
                <option value={ImapOperation.IMAP_OP_FETCH_MESSAGE}>fetch message</option>
            </select>
        </div>
    )
}

function ConnectionPicker({connections, selectedId, onSelect}: {
    connections: ExternalConnectionInfo[]
    selectedId: string
    onSelect: (id: string) => void
}) {
    const navigate = useNavigate()

    if (connections.length === 0) {
        return (
            <div className={cls.ParamsSection}>
                <span className={cls.ParamsSectionTitle}>Connection</span>
                <Button
                    variant="ghost"
                    onClick={() => navigate(Path.ConnectionsPage)}>
                    No connections — add one
                </Button>
            </div>
        )
    }

    return (
        <div className={cls.ParamsSection}>
            <span className={cls.ParamsSectionTitle}>Connection</span>
            {connections.map(c => (
                <SelectOption
                    key={c.id}
                    label={connectionLabel(c)}
                    selected={c.id === selectedId}
                    onSelect={() => onSelect(c.id ?? "")}
                />
            ))}
        </div>
    )
}

function ParamsList({params, values, onChange}: {
    params: Record<string, ToolParamDef>
    values: Record<string, string>
    onChange: (name: string, value: string) => void
}) {
    const entries = Object.entries(params)
    if (entries.length === 0) return null
    return (
        <div className={cls.ParamsSection}>
            <span className={cls.ParamsSectionTitle}>Parameters</span>
            {entries.map(([name, def]) => (
                <ParamRow key={name} name={name} def={def} value={values[name] ?? ""} onChange={onChange}/>
            ))}
        </div>
    )
}

function ParamRow({name, def, value, onChange}: {
    name: string
    def: ToolParamDef
    value: string
    onChange: (name: string, value: string) => void
}) {
    return (
        <div className={cls.ParamRow}>
            <span className={cls.ParamName}>{name}</span>
            {def.description && <span className={cls.ParamDesc}>{def.description}</span>}
            <ParamInput def={def} value={value} onChange={v => onChange(name, v)}/>
        </div>
    )
}

function ParamInput({def, value, onChange}: {
    def: ToolParamDef
    value: string
    onChange: (value: string) => void
}) {
    if (def.enumParam) {
        return (
            <select className={cls.ParamField} value={value} onChange={e => onChange(e.target.value)}>
                <option value="">—</option>
                {(def.enumParam.values ?? []).map(v => (
                    <option key={v} value={v}>{v}</option>
                ))}
            </select>
        )
    }
    if (def.integerParam) {
        return <input className={cls.ParamField} type="number" value={value} onChange={e => onChange(e.target.value)} placeholder="integer"/>
    }
    return <input className={cls.ParamField} type="text" value={value} onChange={e => onChange(e.target.value)} placeholder="string"/>
}

function ResultPanel({result}: { result: string }) {
    return (
        <div className={cls.ParamsSection}>
            <span className={cls.ParamsSectionTitle}>Result</span>
            <pre className={cls.ResultPanel}>{result}</pre>
        </div>
    )
}
