import cls from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/components/ResultPanel/ResultPanel.module.css"

export default function ResultPanel({result}: { result: string }) {
    return (
        <div className={cls.ResultPanelContainer}>
            <span className={cls.SectionTitle}>Result</span>
            <pre className={cls.ResultPre}>{result}</pre>
        </div>
    )
}
