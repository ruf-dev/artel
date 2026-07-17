import {Button} from "@vervstack/chures"

import {cn} from "@/app/utils/cn.ts"
import cls
    // eslint-disable-next-line max-len -- deep nested import path can't be shortened without changing the import
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/segments/ResultView/components/ViewModeToggle/ViewModeToggle.module.css"

export type ResultViewMode = "widget" | "json"

export default function ViewModeToggle({mode, widgetLabel, onChange}: {
    mode: ResultViewMode; widgetLabel: string; onChange: (mode: ResultViewMode) => void
}) {
    return (
        <div className={cls.ViewModeToggleContainer}>
            <Button
                variant="unstyled"
                className={cn(cls.Option, mode === "widget" && cls.OptionActive)}
                onClick={() => onChange("widget")}
            >
                {widgetLabel}
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
