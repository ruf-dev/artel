import {
    LLM_CALL_DESC, LLM_CALL_NAME, LogicOption,
} from "@/pages/tract-canvas/components/TractBlockPicker/useTractBlockPickerData.ts"
import LogicCell from "@/pages/tract-canvas/components/TractBlockPicker/components/LogicCell/LogicCell.tsx"
import OptionCell from "@/pages/tract-canvas/components/TractBlockPicker/components/OptionCell/OptionCell.tsx"
import {ChatIcon} from "@/pages/tract-canvas/icons/ChatIcon/ChatIcon.tsx"
import {colorForKind} from "@/pages/tract-canvas/icons/tractIconHelpers.ts"
import cls from "@/pages/tract-canvas/components/TractBlockPicker/TractBlockPicker.module.css"

interface Props {
    logicOptions: LogicOption[]
    showLlmCall: boolean
    onSelectLogic: (type: LogicOption["type"]) => void
    onSelectLlmCall: () => void
}

/** LogicSection renders the "Logic" category grid (condition/script plus the "Call LLM" option)
 * — split out of TractBlockPicker purely to keep that file under the project's
 * max-lines-per-function limit. */
export default function LogicSection({logicOptions, showLlmCall, onSelectLogic, onSelectLlmCall}: Props) {
    if (logicOptions.length === 0 && !showLlmCall) return null

    return (
        <>
            <div className={cls.CatTitle}>Logic</div>
            <div className={cls.Grid}>
                {logicOptions.map(o => (
                    <LogicCell key={o.type} option={o} onSelect={() => onSelectLogic(o.type)}/>
                ))}
                {showLlmCall && (
                    <OptionCell
                        Icon={ChatIcon}
                        color={colorForKind("llm_call")}
                        name={LLM_CALL_NAME}
                        fn={LLM_CALL_DESC}
                        onSelect={onSelectLlmCall}
                    />
                )}
            </div>
        </>
    )
}
