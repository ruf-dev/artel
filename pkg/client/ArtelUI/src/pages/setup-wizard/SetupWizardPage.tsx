import {useEffect, useState, type ReactNode} from "react"
import {useNavigate} from "react-router-dom"

import {Path} from "@/app/routing/Router.tsx"
import {apiPrefix} from "@/app/api/api.ts"
import {SetupWizardAPI, RegistrationMode} from "@/app/api/artel"
import {
    SetupWizardContext,
    type SetupWizardState,
    type SetupWizardStep,
} from "@/pages/setup-wizard/setupWizardContext.ts"
import TokenEntryScreen from "@/pages/setup-wizard/screens/TokenEntryScreen.tsx"
import AuthMethodsScreen from "@/pages/setup-wizard/screens/AuthMethodsScreen.tsx"
import RegistrationModeScreen from "@/pages/setup-wizard/screens/RegistrationModeScreen.tsx"
import CreateAdminScreen from "@/pages/setup-wizard/screens/CreateAdminScreen.tsx"

export default function SetupWizardPage() {
    const navigate = useNavigate()

    const [step, setStep] = useState<SetupWizardStep>("token")
    const [wizardSessionToken, setWizardSessionToken] = useState("")
    const [passwordEnabled, setPasswordEnabled] = useState(true)
    const [telegramEnabled, setTelegramEnabled] = useState(true)
    const [registrationMode, setRegistrationMode] = useState<RegistrationMode>(RegistrationMode.ADMIN_ONLY)

    // Robust check against the server rather than the possibly-stale
    // useAppConfig value: prevents revisiting a finished wizard (e.g. someone
    // completes setup in another tab, then navigates back to /setup here).
    useEffect(() => {
        SetupWizardAPI.GetStatus({}, apiPrefix())
            .then((res) => {
                if (res.setupCompleted === true) {
                    navigate(Path.HomePage, {replace: true})
                }
            })
            .catch(() => {})
    }, [])

    const value: SetupWizardState = {
        step, setStep,
        wizardSessionToken, setWizardSessionToken,
        passwordEnabled, setPasswordEnabled,
        telegramEnabled, setTelegramEnabled,
        registrationMode, setRegistrationMode,
    }

    return (
        <SetupWizardContext.Provider value={value}>
            {STEP_SCREENS[step]}
        </SetupWizardContext.Provider>
    )
}

const STEP_SCREENS: Record<SetupWizardStep, ReactNode> = {
    token: <TokenEntryScreen/>,
    methods: <AuthMethodsScreen/>,
    mode: <RegistrationModeScreen/>,
    admin: <CreateAdminScreen/>,
}
