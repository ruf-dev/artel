import {create} from "zustand"

interface AppConfigState {
    noAuthEnabled: boolean
    setNoAuthEnabled: (v: boolean) => void
    credsEncrypted: boolean
    setCredsEncrypted: (v: boolean) => void
}

export const useAppConfig = create<AppConfigState>((set) => ({
    noAuthEnabled: false,
    setNoAuthEnabled: (v: boolean) => set({noAuthEnabled: v}),
    // Default true so nothing flashes an insecure warning before the first
    // successful GetConfig response comes back.
    credsEncrypted: true,
    setCredsEncrypted: (v: boolean) => set({credsEncrypted: v}),
}))
