import {Button, Toggle} from "@vervstack/chures"

import cls from "@/pages/setup-wizard/SetupWizardPage.module.css"
import {cn} from "@/app/utils/cn.ts"
import {useSetupWizard} from "@/pages/setup-wizard/setupWizardContext.ts"

export default function AuthMethodsScreen() {
    const {
        passwordEnabled, setPasswordEnabled,
        telegramEnabled, setTelegramEnabled,
        setStep,
    } = useSetupWizard()

    const atLeastOneEnabled = passwordEnabled || telegramEnabled

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
                    onClick={() => setStep("mode")}
                    disabled={!atLeastOneEnabled}
                >
                    Continue
                </Button>
            </div>
        </div>
    )
}
