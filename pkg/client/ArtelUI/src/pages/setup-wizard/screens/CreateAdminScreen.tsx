import {useEffect, useState} from "react"
import {useNavigate} from "react-router-dom"
import {Button, TelegramAuth, type TelegramAuthData} from "@vervstack/chures"

import Input from "@/components/atoms/Input/Input.tsx"
import cls from "@/pages/setup-wizard/SetupWizardPage.module.css"
import {Path} from "@/app/routing/Router.tsx"
import {apiPrefix} from "@/app/api/api.ts"
import {AuthAPI, SetupWizardAPI, type CompleteSetupRequest} from "@/app/api/artel"
import {authService, type Session} from "@/processes/Auth.ts"
import useUser from "@/hooks/user/User.ts"
import {useBakeError} from "@/app/hooks/useErrorToast"
import {isGrpcError, grpcErrorMessage} from "@/processes/grpcErrors"
import {useSetupWizard} from "@/pages/setup-wizard/setupWizardContext.ts"

function isSetupAlreadyCompleted(err: unknown): boolean {
    return isGrpcError(err) && grpcErrorMessage(err).toLowerCase().includes("already been completed")
}

export default function CreateAdminScreen() {
    const navigate = useNavigate()
    const {wizardSessionToken, passwordEnabled, telegramEnabled} = useSetupWizard()
    const bakeError = useBakeError()

    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const [saving, setSaving] = useState(false)
    const [botId, setBotId] = useState("")

    useEffect(() => {
        if (!telegramEnabled) return
        AuthAPI.GetConfig({}, apiPrefix()).then(function (cfg) {
            setBotId(cfg.telegramClientId ?? "")
        }).catch(() => {})
    }, [telegramEnabled])

    function completeSetup(req: CompleteSetupRequest) {
        return SetupWizardAPI.CompleteSetup(req, apiPrefix())
            .then((res) => {
                const session: Session = {token: res.token ?? ""}
                const {login} = useUser.getState()
                login(session)
                return authService.FetchUserInfo()
                    .then((userInfo) => login(session, userInfo))
                    .catch(() => {
                        // profile unavailable — proceed without admin flag
                    })
            })
            .then(() => {
                navigate(Path.HomePage)
            })
            .catch((err: unknown) => {
                if (isSetupAlreadyCompleted(err)) {
                    bakeError("Setup already completed", err)
                    navigate(Path.InitPage)
                    return
                }
                bakeError("Failed to complete setup", err)
            })
    }

    function handlePasswordSubmit() {
        if (!email || !password || saving) return
        setSaving(true)
        completeSetup({wizardSessionToken, password: {email, password}})
            .finally(() => {
                setSaving(false)
            })
    }

    function handleTelegramSuccess(data: TelegramAuthData) {
        void completeSetup({wizardSessionToken, telegram: {idToken: data.id_token}})
    }

    return (
        <div className={cls.SetupWizardPageContainer}>
            <div className={cls.Card}>
                <h2 className={cls.Title}>Create the admin account</h2>
                <p className={cls.Subtitle}>Whoever completes this step becomes the instance administrator.</p>
                {passwordEnabled && (
                    <>
                        <Input
                            type="email"
                            label="Email"
                            value={email}
                            setValue={setEmail}
                            disabled={saving}
                            autoComplete="email"
                        />
                        <Input
                            type="password"
                            label="Password"
                            value={password}
                            setValue={setPassword}
                            disabled={saving}
                            autoComplete="new-password"
                        />
                        <Button
                            variant="primary"
                            className={cls.SubmitBtn}
                            onClick={handlePasswordSubmit}
                            disabled={saving || !email || !password}
                        >
                            {saving ? "Creating…" : "Create admin account"}
                        </Button>
                    </>
                )}
                {passwordEnabled && telegramEnabled && (
                    <div className={cls.Divider}>or</div>
                )}
                {telegramEnabled && (
                    <div className={cls.TelegramSection}>
                        <TelegramAuth
                            title="Sign in with Telegram"
                            botId={botId}
                            onSuccess={handleTelegramSuccess}
                        />
                    </div>
                )}
            </div>
        </div>
    )
}
