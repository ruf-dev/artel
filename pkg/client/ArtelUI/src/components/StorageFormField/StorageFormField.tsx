import Input from "@/components/atoms/Input/Input.tsx"
import cls from "@/components/StorageFormField/StorageFormField.module.css"

interface StorageFormFieldProps {
    label: string
    type?: "text" | "password"
    placeholder?: string
    value: string
    onChange: (v: string) => void
    disabled?: boolean
    autoComplete?: string
}

export default function StorageFormField(props: StorageFormFieldProps) {
    return (
        <label className={cls.StorageFormFieldContainer}>
            <span className={cls.FieldLabel}>{props.label}</span>
            <Input
                type={props.type ?? "text"}
                placeholder={props.placeholder}
                value={props.value}
                setValue={props.onChange}
                disabled={props.disabled}
                autoComplete={props.autoComplete}
            />
        </label>
    )
}
