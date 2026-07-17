import {Button} from "@vervstack/chures"

import {cn} from "@/app/utils/cn.ts"
import cls
    // eslint-disable-next-line max-len -- deep nested import path can't be shortened without changing the import
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/segments/ResultView/components/ViewModeToggle/ViewModeToggle.module.css"

export type ResultViewMode = "table" | "json"

export default function ViewModeToggle({mode, onChange}: {
    mode: ResultViewMode; onChange: (mode: ResultViewMode) => void
}) {
    return (
        <div className={cls.ViewModeToggleContainer}>
            <Button
                variant="unstyled"
                className={cn(cls.Option, mode === "table" && cls.OptionActive)}
                onClick={() => onChange("table")}
            >
                Table
            </Button>
            <Button
                variant="unstyled"
                className={cn(cls.Option, mode === "json" && cls.OptionActive)}
                onClick={() => onChange("json")}
            >
                JSON
            </Button>
        </div>
    )
}
