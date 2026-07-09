import {useState} from "react"

import Input from "@/components/shared/Input/Input.tsx"

interface Props {
    label: string
    placeholder: string
    defaultValue?: string
    onChange: (v: string) => void
    disabled?: boolean
    maxLength?: number
    fieldClassName?: string
    labelClassName?: string
    inputClassName?: string
}

export default function FormField(props: Props) {
    const [value, setValue] = useState(props.defaultValue || "");

    function onValueChange(v: string) {
        setValue(v)
        props.onChange(v)
    }

    return (
        <label className={props.fieldClassName}>
            <span className={props.labelClassName}>{props.label}</span>
            <Input
                className={props.inputClassName}
                placeholder={props.placeholder}
                value={value}
                onChange={e => onValueChange(e.target.value)}
                disabled={props.disabled}
                maxLength={props.maxLength}
                autoComplete="off"
            />
        </label>
    )
}
