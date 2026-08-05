import {create} from "zustand"

interface AppConfigState {
    noAuthEnabled: boolean
    setNoAuthEnabled: (v: boolean) => void
    credsEncrypted: boolean
    setCredsEncrypted: (v: boolean) => void
    setupCompleted: boolean
    setSetupCompleted: (v: boolean) => void
    passwordAuthEnabled: boolean
    setPasswordAuthEnabled: (v: boolean) => void
    telegramAuthEnabled: boolean
    setTelegramAuthEnabled: (v: boolean) => void
    selfRegisterEnabled: boolean
    setSelfRegisterEnabled: (v: boolean) => void
}

export const useAppConfig = create<AppConfigState>((set) => ({
    noAuthEnabled: false,
    setNoAuthEnabled: (v: boolean) => set({noAuthEnabled: v}),
    // Default true so nothing flashes an insecure warning before the first
    // successful GetConfig response comes back.
    credsEncrypted: true,
    setCredsEncrypted: (v: boolean) => set({credsEncrypted: v}),
    // Default true so nothing flashes a setup-wizard redirect before the
    // first successful GetConfig response comes back.
    setupCompleted: true,
    setSetupCompleted: (v: boolean) => set({setupCompleted: v}),
    // Safe defaults: assume both auth methods are available and
    // registration is closed until GetConfig says otherwise.
    passwordAuthEnabled: true,
    setPasswordAuthEnabled: (v: boolean) => set({passwordAuthEnabled: v}),
    telegramAuthEnabled: true,
    setTelegramAuthEnabled: (v: boolean) => set({telegramAuthEnabled: v}),
    selfRegisterEnabled: false,
    setSelfRegisterEnabled: (v: boolean) => set({selfRegisterEnabled: v}),
}))
