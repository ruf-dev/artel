import {Button} from "@vervstack/chures"

import cls from "@/pages/setup-wizard/SetupWizardPage.module.css"
import {cn} from "@/app/utils/cn.ts"
import {apiPrefix} from "@/app/api/api.ts"
import {SetupWizardAPI, RegistrationMode} from "@/app/api/artel"
import {useBakeError} from "@/app/hooks/useErrorToast"
import {useSetupWizard} from "@/pages/setup-wizard/setupWizardContext.ts"

const OPTIONS: { mode: RegistrationMode; title: string; description: string }[] = [
    {
        mode: RegistrationMode.ADMIN_ONLY,
        title: "Admin-only",
        description: "Only admins create accounts",
    },
    {
        mode: RegistrationMode.SELF_REGISTER,
        title: "Open registration",
        description: "Anyone can register",
    },
]

export default function RegistrationModeScreen() {
    const {wizardSessionToken, registrationMode, setRegistrationMode, setStep} = useSetupWizard()
    const bakeError = useBakeError()

    function handleContinue() {
        SetupWizardAPI.SelectRegistrationMode({wizardSessionToken, mode: registrationMode}, apiPrefix())
            .then(() => {
                setStep("admin")
            })
            .catch((err: unknown) => {
                bakeError("Failed to save registration mode", err)
            })
    }

    return (
        <div className={cls.SetupWizardPageContainer}>
            <div className={cls.Card}>
                <h2 className={cls.Title}>Registration mode</h2>
                <p className={cls.Subtitle}>Choose how new accounts get created on this instance.</p>
                {OPTIONS.map((opt) => (
                    <Button
                        key={opt.mode}
                        variant="ghost"
                        className={cn(cls.OptionRow, registrationMode === opt.mode && cls.OptionRowActive)}
                        onClick={() => setRegistrationMode(opt.mode)}
                    >
                        <span className={cls.OptionTitle}>{opt.title}</span>
                        <span className={cls.OptionDescription}>{opt.description}</span>
                    </Button>
                ))}
                <Button variant="primary" className={cls.SubmitBtn} onClick={handleContinue}>
                    Continue
                </Button>
            </div>
        </div>
    )
}
