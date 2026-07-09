import cls from "@/components/S3InstanceFormDialog/S3InstanceFormDialog.module.css"
import Input from "@/components/shared/Input/Input.tsx"

interface Props {
    useSsl: boolean
    pathStyle: boolean
    disabled: boolean
    onUseSslChange: (v: boolean) => void
    onPathStyleChange: (v: boolean) => void
}

export default function S3ToggleFields({useSsl, pathStyle, disabled, onUseSslChange, onPathStyleChange}: Props) {
    return (
        <div className={cls.ToggleGroup}>
            <label className={cls.ToggleRow}>
                <Input
                    type="checkbox"
                    checked={useSsl}
                    disabled={disabled}
                    onChange={e => onUseSslChange(e.currentTarget.checked)}
                />
                <span className={cls.ToggleLabel}>Use SSL</span>
            </label>
            <label className={cls.ToggleRow}>
                <Input
                    type="checkbox"
                    checked={pathStyle}
                    disabled={disabled}
                    onChange={e => onPathStyleChange(e.currentTarget.checked)}
                />
                <span className={cls.ToggleLabel}>Path-style addressing</span>
            </label>
            <p className={cls.ToggleCaption}>
                Required for Garage and most self-hosted S3-compatible stores.
            </p>
        </div>
    )
}
