import {Button} from "@vervstack/chures"

import cls from "@/dialogs/ConnectGenericDialog/components/CredentialRow/CredentialRow.module.css"
import Input from "@/components/atoms/Input/Input.tsx"

interface CredentialRowProps {
    keyValue: string
    value: string
    onKeyChange: (v: string) => void
    onValueChange: (v: string) => void
    onRemove: () => void
    disabled: boolean
    canRemove: boolean
}

export default function CredentialRow(props: CredentialRowProps) {
    return (
        <div className={cls.CredentialRowContainer}>
            <Input
                className={cls.KeyInput} placeholder="Field name (e.g. api_key)" value={props.keyValue}
                setValue={props.onKeyChange} disabled={props.disabled} autoComplete="off"
            />
            <Input
                className={cls.ValueInput} placeholder="Value" value={props.value}
                setValue={props.onValueChange} disabled={props.disabled} autoComplete="off"
            />
            {props.canRemove && (
                <Button
                    variant="ghost" className={cls.RemoveBtn} onClick={props.onRemove} disabled={props.disabled}
                    aria-label="Remove field"
                >
                    ✕
                </Button>
            )}
        </div>
    )
}
