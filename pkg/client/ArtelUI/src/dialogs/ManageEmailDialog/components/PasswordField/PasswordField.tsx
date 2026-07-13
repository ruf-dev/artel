import Input from "@/components/atoms/Input/Input.tsx"
import cls from "@/dialogs/ManageEmailDialog/components/PasswordField/PasswordField.module.css"

export default function PasswordField({value, onChange, disabled, appPasswordUrl, hasStoredValue}: {
    value: string
    onChange: (value: string) => void
    disabled: boolean
    appPasswordUrl: string
    hasStoredValue: boolean
}) {
    const placeholder = hasStoredValue && !value ? "••••••••" : "App-specific password"

    return (
        <label className={cls.PasswordFieldContainer}>
            <span className={cls.FieldLabel}>Password</span>
            <Input type="password" placeholder={placeholder} value={value}
                setValue={onChange} disabled={disabled} autoComplete="new-password"/>
            {appPasswordUrl && (
                <div className={cls.YandexHint}>
                    <a href={appPasswordUrl} target="_blank" rel="noopener noreferrer">
                        Click here to create new password
                    </a>
                </div>
            )}
        </label>
    )
}
