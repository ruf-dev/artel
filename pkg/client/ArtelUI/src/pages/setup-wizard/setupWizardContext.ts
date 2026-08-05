import {createContext, useContext, type Dispatch, type SetStateAction} from "react"

import {RegistrationMode} from "@/app/api/artel"

export type SetupWizardStep = "token" | "methods" | "mode" | "admin"

export interface SetupWizardState {
    step: SetupWizardStep
    setStep: Dispatch<SetStateAction<SetupWizardStep>>

    wizardSessionToken: string
    setWizardSessionToken: Dispatch<SetStateAction<string>>

    passwordEnabled: boolean
    setPasswordEnabled: Dispatch<SetStateAction<boolean>>
    telegramEnabled: boolean
    setTelegramEnabled: Dispatch<SetStateAction<boolean>>

    registrationMode: RegistrationMode
    setRegistrationMode: Dispatch<SetStateAction<RegistrationMode>>
}

export const SetupWizardContext = createContext<SetupWizardState | null>(null)

export function useSetupWizard(): SetupWizardState {
    const ctx = useContext(SetupWizardContext)
    if (!ctx) throw new Error("useSetupWizard must be used within SetupWizardContext.Provider")
    return ctx
}
