import {create} from 'zustand'

import {AuthMiddleware, UserInfo} from "@/processes/AuthMiddleware.ts"
import { Session} from "@/processes/Auth.ts"

interface UserState {
    auth: AuthMiddleware
    isAdmin: boolean
    isEmailsEnabled: boolean

    login: (session: Session, userInfo?: UserInfo) => void
    logout: () => void
    setUserInfo: (info: UserInfo) => void
}

const useUser = create<UserState>((set, get) => ({
    auth: new AuthMiddleware(),
    isAdmin: new AuthMiddleware().isAdmin(),
    isEmailsEnabled: new AuthMiddleware().hasEmailsPermission(),

    login: (session: Session, userInfo?: UserInfo) => {
        get().auth.login(session, userInfo)
        set({auth: get().auth, isAdmin: get().auth.isAdmin(), isEmailsEnabled: get().auth.hasEmailsPermission()})
    },

    logout: () => {
        get().auth.logout()
        set({auth: new AuthMiddleware(), isAdmin: false, isEmailsEnabled: false})
    },

    setUserInfo: (info: UserInfo) => {
        get().auth.setUserInfo(info)
        set({auth: get().auth, isAdmin: get().auth.isAdmin(), isEmailsEnabled: get().auth.hasEmailsPermission()})
    },
}))

export default useUser
