import cls from "@/pages/tract-canvas/components/OutputFields/OutputFields.module.css"

import {SchemaNode} from "@/processes/Tracts.ts"
import {cn} from "@/app/utils/cn.ts"

interface Props {
    schema?: SchemaNode
}

export default function OutputFields({schema}: Props) {
    const props = schema ? Object.entries(schema.properties) : []

    return (
        <div className={cls.OutputFieldsContainer}>
            {props.length === 0 ? (
                <p className={cls.Empty}>No output fields.</p>
            ) : (
                props.map(([name, def]) => (
                    <div className={cls.Row} key={name}>
                        <span className={cls.Key}>{name}</span>
                        <span className={cn(cls.Val, cls.ValDim)}>{def.type}{def.description ? ` — ${def.description}` : ""}</span>
                    </div>
                ))
            )}
        </div>
    )
}
