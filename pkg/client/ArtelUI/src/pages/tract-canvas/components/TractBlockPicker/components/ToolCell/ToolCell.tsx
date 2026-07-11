import {TractTool} from "@/processes/Tracts.ts"
import {colorForKind, iconForTool} from "@/pages/tract-canvas/icons/tractIconHelpers.ts"
import OptionCell from "@/pages/tract-canvas/components/TractBlockPicker/components/OptionCell/OptionCell.tsx"

interface ToolCellProps {
    tool: TractTool
    onSelect: () => void
}

export default function ToolCell({tool, onSelect}: ToolCellProps) {
    const color = colorForKind("action", tool.mcp)
    const Icon = iconForTool(tool.mcp, tool.tool)
    return <OptionCell Icon={Icon} color={color} name={tool.tool} fn={`${tool.mcp}.${tool.tool}`} onSelect={onSelect}/>
}
