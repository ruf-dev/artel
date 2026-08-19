import {useState} from "react"
import {Button, Input} from "@vervstack/chures"

import cls from "@/pages/workbench/components/WorkbenchAuthScreen/components/AuthCodeForm/AuthCodeForm.module.css"

interface Props {
    onSubmit: (code: string) => void
    disabled?: boolean
}

export default function AuthCodeForm({onSubmit, disabled}: Props) {
    const [code, setCode] = useState("")

    function handleSubmit() {
        const trimmed = code.trim()
        if (!trimmed || disabled) return
        onSubmit(trimmed)
        setCode("")
    }

    function handleKeyDown(e: React.KeyboardEvent<HTMLInputElement>) {
        if (e.key === "Enter" && code.trim() && !disabled) {
            handleSubmit()
        }
    }

    return (
        <div className={cls.AuthCodeFormContainer}>
            <Input
                value={code}
                setValue={setCode}
                onKeyDown={handleKeyDown}
                placeholder="Enter the authentication code…"
                disabled={disabled}
                className={cls.InputWrapper}
            />
            <Button variant="primary" onClick={handleSubmit} disabled={disabled || !code.trim()}>
                Submit
            </Button>
        </div>
    )
}
