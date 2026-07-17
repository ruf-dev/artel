import JsonView from "@/components/JsonView/JsonView.tsx"
import cls
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/components/ToolResponseViewer/ToolResponseViewer.module.css"

function tryParseJson(result: string): {ok: true; value: unknown} | {ok: false} {
    try {
        return {ok: true, value: JSON.parse(result)}
    } catch {
        return {ok: false}
    }
}

export default function ToolResponseViewer({result}: { result: string }) {
    const parsed = tryParseJson(result)
    return (
        <div className={cls.ToolResponseViewerContainer}>
            <span className={cls.SectionTitle}>Result</span>
            {parsed.ok ? <JsonView value={parsed.value}/> : <pre className={cls.ResponsePre}>{result}</pre>}
        </div>
    )
}
