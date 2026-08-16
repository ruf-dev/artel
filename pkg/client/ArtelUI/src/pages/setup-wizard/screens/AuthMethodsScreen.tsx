import {Button, Toggle} from "@vervstack/chures"

import cls from "@/pages/setup-wizard/SetupWizardPage.module.css"
import {cn} from "@/app/utils/cn.ts"
import {apiPrefix} from "@/app/api/api.ts"
import {SetupWizardAPI} from "@/app/api/artel"
import {useBakeError} from "@/app/hooks/useErrorToast"
import {useSetupWizard} from "@/pages/setup-wizard/setupWizardContext.ts"

export default function AuthMethodsScreen() {
    const {
        wizardSessionToken,
        passwordEnabled, setPasswordEnabled,
        telegramEnabled, setTelegramEnabled,
        setStep,
    } = useSetupWizard()
    const bakeError = useBakeError()

    const atLeastOneEnabled = passwordEnabled || telegramEnabled

    function handleContinue() {
        SetupWizardAPI.SelectAuthMethods({wizardSessionToken, passwordEnabled, telegramEnabled}, apiPrefix())
            .then(() => {
                setStep("mode")
            })
            .catch((err: unknown) => {
                bakeError("Failed to save sign-in methods", err)
            })
    }

    return (
        <div className={cls.SetupWizardPageContainer}>
            <div className={cls.Card}>
                <h2 className={cls.Title}>Choose sign-in methods</h2>
                <p className={cls.Subtitle}>Pick which ways people can sign in to this instance.</p>
                <Toggle
                    checked={passwordEnabled}
                    onChange={setPasswordEnabled}
                    label="Password login"
                    labelPosition="left"
                />
                <Toggle
                    checked={telegramEnabled}
                    onChange={setTelegramEnabled}
                    label="Telegram login"
                    labelPosition="left"
                />
                {!atLeastOneEnabled && (
                    <p className={cn(cls.Hint, cls.HintError)}>At least one sign-in method must stay enabled.</p>
                )}
                <Button
                    variant="primary"
                    className={cls.SubmitBtn}
                    onClick={handleContinue}
                    disabled={!atLeastOneEnabled}
                >
                    Continue
                </Button>
            </div>
        </div>
    )
}
