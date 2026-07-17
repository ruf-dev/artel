import cls
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/components/ToolResponseViewer/ToolResponseViewer.module.css"

export default function ToolResponseViewer({result}: { result: string }) {
    return (
        <div className={cls.ToolResponseViewerContainer}>
            <span className={cls.SectionTitle}>Result</span>
            <pre className={cls.ResponsePre}>{result}</pre>
        </div>
    )
}
