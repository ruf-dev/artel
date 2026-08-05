import {useState} from "react"
import {Button} from "@vervstack/chures"

import Input from "@/components/atoms/Input/Input.tsx"
import cls from "@/pages/setup-wizard/SetupWizardPage.module.css"
import {apiPrefix} from "@/app/api/api.ts"
import {SetupWizardAPI} from "@/app/api/artel"
import {useBakeError} from "@/app/hooks/useErrorToast"
import {useSetupWizard} from "@/pages/setup-wizard/setupWizardContext.ts"

export default function TokenEntryScreen() {
    const {setStep, setWizardSessionToken} = useSetupWizard()
    const bakeError = useBakeError()

    const [token, setToken] = useState("")
    const [submitting, setSubmitting] = useState(false)

    function handleSubmit() {
        if (!token || submitting) return
        setSubmitting(true)
        SetupWizardAPI.SubmitToken({token}, apiPrefix())
            .then((res) => {
                setWizardSessionToken(res.wizardSessionToken ?? "")
                setStep("methods")
            })
            .catch((err: unknown) => {
                bakeError("Invalid setup token", err)
            })
            .finally(() => {
                setSubmitting(false)
            })
    }

    return (
        <div className={cls.SetupWizardPageContainer}>
            <div className={cls.Card}>
                <div className={cls.Logo}>artel</div>
                <h2 className={cls.Title}>Welcome — let's set up your instance</h2>
                <p className={cls.Subtitle}>Enter the one-time setup token to begin.</p>
                <Input
                    type="password"
                    label="Setup token"
                    value={token}
                    setValue={setToken}
                    disabled={submitting}
                    autoComplete="off"
                    onKeyDown={(e) => {
                        if (e.key === "Enter") handleSubmit()
                    }}
                />
                <p className={cls.Hint}>Find this in your server logs — grep for "setup token:"</p>
                <Button
                    variant="primary"
                    className={cls.SubmitBtn}
                    onClick={handleSubmit}
                    disabled={submitting || !token}
                >
                    {submitting ? "Checking…" : "Continue"}
                </Button>
            </div>
        </div>
    )
}
