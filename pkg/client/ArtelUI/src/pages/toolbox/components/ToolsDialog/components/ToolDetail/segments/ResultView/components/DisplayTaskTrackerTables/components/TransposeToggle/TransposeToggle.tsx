import {Button} from "@vervstack/chures"

import {cn} from "@/app/utils/cn.ts"
import cls
    // eslint-disable-next-line max-len -- deep nested import path can't be shortened without changing the import
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/segments/ResultView/components/DisplayTaskTrackerTables/components/TransposeToggle/TransposeToggle.module.css"

export default function TransposeToggle({transposed, onChange}: {
    transposed: boolean; onChange: (transposed: boolean) => void
}) {
    return (
        <div className={cls.TransposeToggleContainer}>
            <Button
                variant="unstyled"
                className={cn(cls.ToggleBtn, transposed && cls.ToggleBtnActive)}
                onClick={() => onChange(!transposed)}
            >
                Transpose
            </Button>
        </div>
    )
}
