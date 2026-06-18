import {useEffect, useState} from "react"
import {useNavigate} from "react-router-dom"

import cls from "@/pages/toolbox/ToolboxPage.module.css"

import {ImapOperation, ImapToolAction, McpToolInfo, MomCandidate, SmtpOperation, SmtpToolAction, ToolParamDef} from "@/app/api/artel/mcp_keys.pb.ts"
import {Path} from "@/app/routing/Router.tsx"
import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import useUser from "@/hooks/user/User.ts"
import {useDialog} from "@/app/hooks/Dialog"
import Button from "@/components/shared/Button/Button.tsx"

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
        return (
            <div className={cls.DialogContainer} role="dialog" aria-modal="true">
                <h2 className={cls.DialogTitle}>{activeTool.name}</h2>
                <p className={cls.DialogDesc}>{activeTool.description}</p>
                <div className={cls.ToolDetail}>
                    {activeTool.smtp && <SmtpActionView action={activeTool.smtp}/>}
                    {activeTool.imap && <ImapActionView action={activeTool.imap}/>}
                    {activeTool.params && <ParamsList params={activeTool.params}/>}
                </div>
                <div className={cls.DialogActions}>
                    <Button variant="ghost" onClick={() => setActiveTool(null)}>Back</Button>
                </div>
            </div>
        )
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

function ParamsList({params}: { params: Record<string, ToolParamDef> }) {
    const entries = Object.entries(params)
    if (entries.length === 0) return null
    return (
        <div className={cls.ParamsSection}>
            <span className={cls.ParamsSectionTitle}>Parameters</span>
            {entries.map(([name, def]) => (
                <ParamRow key={name} name={name} def={def}/>
            ))}
        </div>
    )
}

function ParamRow({name, def}: { name: string; def: ToolParamDef }) {
    return (
        <div className={cls.ParamRow}>
            <span className={cls.ParamName}>{name}</span>
            {def.description && <span className={cls.ParamDesc}>{def.description}</span>}
            <ParamInput def={def}/>
        </div>
    )
}

function ParamInput({def}: { def: ToolParamDef }) {
    if (def.enumParam) {
        return (
            <select className={cls.ActionSelect} disabled>
                {(def.enumParam.values ?? []).map(v => (
                    <option key={v} value={v}>{v}</option>
                ))}
            </select>
        )
    }
    if (def.integerParam) {
        return <input className={cls.ParamField} type="number" disabled placeholder="integer"/>
    }
    return <input className={cls.ParamField} type="text" disabled placeholder="string"/>
}
