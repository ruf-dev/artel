import {useState} from "react"

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

export default function FormField({
                                      label,
                                      placeholder,
                                      defaultValue,
                                      onChange,
                                      disabled,
                                      maxLength,
                                      fieldClassName,
                                      labelClassName,
                                      inputClassName
                                  }: Props) {

    const [value, setValue] = useState(defaultValue || "");

    function onValueChange(v: string) {
        setValue(v)
        onChange(v)
    }

    return (
        <label className={fieldClassName}>
            <span className={labelClassName}>{label}</span>
            <input
                className={inputClassName}
                placeholder={placeholder}
                value={value}
                onChange={e => onValueChange(e.target.value)}

                disabled={disabled}
                maxLength={maxLength}
                autoComplete="off"
            />
        </label>
    )
}
