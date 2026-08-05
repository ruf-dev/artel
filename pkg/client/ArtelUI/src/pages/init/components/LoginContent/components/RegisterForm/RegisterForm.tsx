import {useState} from "react"
import {Button} from "@vervstack/chures"

import Input from "@/components/atoms/Input/Input.tsx"
import {authService, Session} from "@/processes/Auth.ts"
import {useBakeError} from "@/app/hooks/useErrorToast"
import cls from "@/pages/init/components/LoginContent/components/RegisterForm/RegisterForm.module.css"

interface RegisterFormProps {
    onRegistered: (session: Session) => void
    onSwitchToLogin: () => void
}

export default function RegisterForm({onRegistered, onSwitchToLogin}: RegisterFormProps) {
    const bakeError = useBakeError()
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const [confirmPassword, setConfirmPassword] = useState("")
    const [submitting, setSubmitting] = useState(false)

    function handleSubmit() {
        if (password !== confirmPassword) {
            bakeError("Registration failed", new Error("Passwords do not match"))
            return
        }
        setSubmitting(true)
        authService.Register(email, password)
            .then(() => authService.LoginWithPassword(email, password))
            .then(onRegistered)
            .catch((err: unknown) => bakeError("Registration failed", err))
            .finally(() => setSubmitting(false))
    }

    return (
        <div className={cls.RegisterFormContainer}>
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
                autoComplete="new-password"
            />
            <Input
                type="password"
                placeholder="Confirm password"
                value={confirmPassword}
                setValue={setConfirmPassword}
                disabled={submitting}
                autoComplete="new-password"
            />
            <Button
                type="button"
                variant="primary"
                className={cls.SubmitButton}
                disabled={submitting || !email || !password || !confirmPassword}
                onClick={handleSubmit}
            >
                {submitting ? "Creating account…" : "Create account"}
            </Button>
            <Button
                type="button"
                variant="unstyled"
                className={cls.SwitchLink}
                onClick={onSwitchToLogin}
            >
                Already have an account? Sign in
            </Button>
        </div>
    )
}
