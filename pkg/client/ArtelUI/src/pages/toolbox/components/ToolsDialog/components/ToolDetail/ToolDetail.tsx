import {useState} from "react"
import {Button} from "@vervstack/chures"

import {McpToolInfo, MomCandidate} from "@/app/api/artel/mcp_keys.pb.ts"
import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import {useBakeError} from "@/app/hooks/useErrorToast"
import {coerceParams} from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/processes/coerceParams.ts"
import SmtpActionView
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/components/SmtpActionView/SmtpActionView.tsx"
import ImapActionView
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/components/ImapActionView/ImapActionView.tsx"
import ConnectionPicker
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/components/ConnectionPicker/ConnectionPicker.tsx"
import ParamsList
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/components/ParamsList/ParamsList.tsx"
import RunScreens
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/components/RunScreens/RunScreens.tsx"
import cls from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/ToolDetail.module.css"

export default function ToolDetail(
    {candidate, tool, onBack}: { candidate: MomCandidate; tool: McpToolInfo; onBack: () => void }
) {
    const {executeMomTool} = useMcpKeys()
    const bakeError = useBakeError()
    const connections = candidate.connections ?? []
    const [selectedConnectionId, setSelectedConnectionId] = useState<string>(
        connections.length === 1 ? (connections[0].id ?? "") : ""
    )
    const [paramValues, setParamValues] = useState<Record<string, string>>({})
    const [running, setRunning] = useState(false)
    const [result, setResult] = useState<string | null>(null)
    const [showResult, setShowResult] = useState(false)

    function setParam(name: string, value: string) {
        setParamValues(prev => ({...prev, [name]: value}))
    }

    function onRun() {
        if (!selectedConnectionId) return
        setRunning(true)
        setShowResult(false)
        const params = coerceParams(tool.params ?? {}, paramValues)
        executeMomTool(candidate.name ?? "", tool.name ?? "", selectedConnectionId, params)
            .then(res => {
                setResult(res)
                setShowResult(true)
            })
            .catch(err => bakeError("Failed to run tool", err))
            .finally(() => setRunning(false))
    }

    function onBackClick() {
        if (showResult) {
            setShowResult(false)
            return
        }
        onBack()
    }

    return (
        <div className={cls.ToolDetailContainer} role="dialog" aria-modal="true">
            <h2 className={cls.DialogTitle}>{tool.name}</h2>
            <p className={cls.DialogDesc}>{tool.description}</p>
            <RunScreens showResult={showResult} result={result} candidate={candidate} tool={tool}>
                {tool.smtp && <SmtpActionView action={tool.smtp}/>}
                {tool.imap && <ImapActionView action={tool.imap}/>}
                <ConnectionPicker
                    connections={connections}
                    selectedId={selectedConnectionId}
                    onSelect={setSelectedConnectionId}
                />
                {tool.params && <ParamsList params={tool.params} values={paramValues} onChange={setParam}/>}
            </RunScreens>
            <div className={cls.DialogActions}>
                <Button variant="ghost" onClick={onBackClick}>Back</Button>
                <Button variant="primary" disabled={running || !selectedConnectionId} onClick={onRun}>
                    {running ? "Running…" : showResult ? "Run again" : "Run"}
                </Button>
            </div>
        </div>
    )
}
