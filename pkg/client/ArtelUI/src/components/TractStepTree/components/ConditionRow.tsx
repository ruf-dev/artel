import {Button} from "@vervstack/chures"

import TemplateInput from "@/components/TemplateInput/TemplateInput.tsx"
import {TemplateSource} from "@/components/TemplateInput/processes/templateRefs.ts"
import {TractCondition} from "@/processes/Tracts.ts"
import cls from "@/components/TractStepTree/TractStepTree.module.css"

const CONDITION_OPS = ["==", "!=", ">", "<", ">=", "<=", "contains", "glob", "regex"] as const

interface ConditionRowProps {
    cond: TractCondition
    sources: TemplateSource[]
    onUpdate: (patch: Partial<TractCondition>) => void
    onRemove: () => void
}

export default function ConditionRow({cond, sources, onUpdate, onRemove}: ConditionRowProps) {
    return (
        <div className={cls.ConditionRow}>
            <TemplateInput
                value={cond.left}
                onChange={v => onUpdate({left: v})}
                sources={sources}
                placeholder="left"
            />
            <select
                className={cls.ConditionOp}
                value={cond.op}
                onChange={e => onUpdate({op: e.target.value as TractCondition["op"]})}
            >
                {CONDITION_OPS.map(op => <option key={op} value={op}>{op}</option>)}
            </select>
            <TemplateInput
                value={cond.right}
                onChange={v => onUpdate({right: v})}
                sources={sources}
                placeholder="right"
            />
            <Button
                variant="iconDanger"
                className={cls.RemoveRowBtn}
                onClick={onRemove}
                aria-label="Remove condition"
            >
                ✕
            </Button>
        </div>
    )
}
