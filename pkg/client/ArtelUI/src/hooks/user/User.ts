import {create} from 'zustand'

import {AuthMiddleware, Session} from "@/processes/Auth.ts"

interface UserState {
    auth: AuthMiddleware
    login: (session: Session) => void
    logout: () => void
}

const useUser = create<UserState>((set, get) => ({
    auth: new AuthMiddleware(),

    login: (session: Session) => {
        get().auth.login(session)
        set({auth: get().auth})
    },

    logout: () => {
        get().auth.logout()
        set({auth: new AuthMiddleware()})
    },
}))

export default useUser
