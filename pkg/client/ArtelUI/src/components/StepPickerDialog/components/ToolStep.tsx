import {Button} from "@vervstack/chures"

import cls from "@/components/StepPickerDialog/StepPickerDialog.module.css"
import {TractTool} from "@/processes/Tracts.ts"
import {useDialog} from "@/app/hooks/Dialog"

interface Props {
    tools: TractTool[]
    onSelect: (tool: TractTool) => void
    onBack: () => void
}

export default function ToolStep({tools, onSelect, onBack}: Props) {
    const {CloseDialog} = useDialog()

    const grouped = tools.reduce<Record<string, TractTool[]>>((acc, t) => {
        (acc[t.mcp] ??= []).push(t)
        return acc
    }, {})

    return (
        <div className={cls.DialogContainer} role="dialog" aria-modal="true">
            <h2 className={cls.DialogTitle}>Choose an action</h2>
            <div className={cls.List}>
                {Object.entries(grouped).map(([mcp, mcpTools]) => (
                    <div className={cls.McpGroup} key={mcp}>
                        <span className={cls.McpGroupLabel}>{mcp}</span>
                        {mcpTools.map(t => (
                            <Button
                                key={t.tool}
                                variant="ghost"
                                className={cls.ToolRow}
                                onClick={() => onSelect(t)}
                            >
                                <span className={cls.ToolName}>{t.tool}</span>
                                {t.description && <span className={cls.ToolDesc}>{t.description}</span>}
                            </Button>
                        ))}
                    </div>
                ))}
                {tools.length === 0 && <p className={cls.Empty}>No tools available.</p>}
            </div>
            <div className={cls.DialogActions}>
                <Button variant="ghost" onClick={onBack}>Back</Button>
                <Button variant="ghost" onClick={CloseDialog}>Cancel</Button>
            </div>
        </div>
    )
}
