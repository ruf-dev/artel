import {useState} from "react"
import {Button} from "@vervstack/chures"

import Input from "@/components/atoms/Input/Input.tsx"
import {authService, Session} from "@/processes/Auth.ts"
import {useBakeError} from "@/app/hooks/useErrorToast"
import cls from "@/pages/init/components/LoginContent/components/PasswordLoginForm/PasswordLoginForm.module.css"

interface PasswordLoginFormProps {
    onLoggedIn: (session: Session) => void
}

export default function PasswordLoginForm({onLoggedIn}: PasswordLoginFormProps) {
    const bakeError = useBakeError()
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const [submitting, setSubmitting] = useState(false)

    function handleSubmit() {
        setSubmitting(true)
        authService.LoginWithPassword(email, password)
            .then(onLoggedIn)
            .catch((err: unknown) => bakeError("Login failed", err))
            .finally(() => setSubmitting(false))
    }

    return (
        <div className={cls.PasswordLoginFormContainer}>
            <Input
                type="email"
                placeholder="Email"
                value={email}
                setValue={setEmail}
                disabled={submitting}
                autoComplete="email"
            />
            <Input
                type="password"
                placeholder="Password"
                value={password}
                setValue={setPassword}
                disabled={submitting}
                autoComplete="current-password"
            />
            <Button
                type="button"
                variant="primary"
                className={cls.SubmitButton}
                disabled={submitting || !email || !password}
                onClick={handleSubmit}
            >
                {submitting ? "Signing in…" : "Sign in"}
            </Button>
        </div>
    )
}
