import Input from "@/components/atoms/Input/Input.tsx"
import cls from "@/dialogs/ManageEmailDialog/components/HostPortRow/HostPortRow.module.css"

interface HostPortRowProps {
    host: { label: string; value: string; placeholder: string; onChange: (v: string) => void }
    port: { label: string; value: string; placeholder: string; onChange: (v: string) => void }
    disabled?: boolean
}

export default function HostPortRow({host, port, disabled}: HostPortRowProps) {
    return (
        <div className={cls.FieldRowContainer}>
            <label className={cls.Field}>
                <span className={cls.FieldLabel}>{host.label}</span>
                <Input placeholder={host.placeholder} value={host.value}
                    setValue={host.onChange} disabled={disabled} autoComplete="off"/>
            </label>
            <label className={cls.FieldNarrow}>
                <span className={cls.FieldLabel}>{port.label}</span>
                <Input placeholder={port.placeholder} value={port.value}
                    setValue={port.onChange} disabled={disabled} maxLength={5}
                    autoComplete="off"/>
            </label>
        </div>
    )
}
