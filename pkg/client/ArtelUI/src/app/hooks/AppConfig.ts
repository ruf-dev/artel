import {create} from "zustand"

interface AppConfigState {
    noAuthEnabled: boolean
    setNoAuthEnabled: (v: boolean) => void
}

export const useAppConfig = create<AppConfigState>((set) => ({
    noAuthEnabled: false,
    setNoAuthEnabled: (v: boolean) => set({noAuthEnabled: v}),
}))
