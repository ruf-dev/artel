import {Button} from "@vervstack/chures"

import {TractTool} from "@/processes/Tracts.ts"
import {colorForKind, iconForTool} from "@/pages/tract-canvas/icons/tractIconHelpers.ts"
import cls from "@/pages/tract-canvas/components/TractBlockPicker/components/ToolCell/ToolCell.module.css"

interface ToolCellProps {
    tool: TractTool
    onSelect: () => void
}

export default function ToolCell({tool, onSelect}: ToolCellProps) {
    const color = colorForKind("action", tool.mcp)
    const Icon = iconForTool(tool.mcp, tool.tool)
    return (
        <Button variant="ghost" className={cls.Opt} onClick={onSelect}>
            <div className={cls.OptIcon} style={{background: color.bg, borderColor: color.border, color: color.fg}}>
                <Icon/>
            </div>
            <div className={cls.OptText}>
                <div className={cls.OptName}>{tool.tool}</div>
                <div className={cls.OptFn}>{tool.mcp}.{tool.tool}</div>
            </div>
        </Button>
    )
}
