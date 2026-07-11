import {create} from 'zustand'

import {AuthMiddleware, UserInfo} from "@/processes/AuthMiddleware.ts"
import { Session} from "@/processes/Auth.ts"

interface UserState {
    auth: AuthMiddleware
    isAdmin: boolean
    isEmailsEnabled: boolean
    isTaskTrackersEnabled: boolean
    isNotesEnabled: boolean
    photoUrl: string | undefined

    login: (session: Session, userInfo?: UserInfo) => void
    logout: () => void
    setUserInfo: (info: UserInfo) => void
}

const useUser = create<UserState>((set, get) => ({
    auth: new AuthMiddleware(),
    isAdmin: new AuthMiddleware().isAdmin(),
    isEmailsEnabled: new AuthMiddleware().hasEmailsPermission(),
    isTaskTrackersEnabled: new AuthMiddleware().hasTaskTrackersPermission(),
    isNotesEnabled: new AuthMiddleware().hasNotesPermission(),
    photoUrl: new AuthMiddleware().getPhotoUrl(),

    login: (session: Session, userInfo?: UserInfo) => {
        get().auth.login(session, userInfo)
        set({
            auth: get().auth,
            isAdmin: get().auth.isAdmin(),
            isEmailsEnabled: get().auth.hasEmailsPermission(),
            isTaskTrackersEnabled: get().auth.hasTaskTrackersPermission(),
            isNotesEnabled: get().auth.hasNotesPermission(),
            photoUrl: get().auth.getPhotoUrl(),
        })
    },

    logout: () => {
        get().auth.logout()
        set({
            auth: new AuthMiddleware(),
            isAdmin: false,
            isEmailsEnabled: false,
            isTaskTrackersEnabled: false,
            isNotesEnabled: false,
            photoUrl: undefined,
        })
    },

    setUserInfo: (info: UserInfo) => {
        get().auth.setUserInfo(info)
        set({
            auth: get().auth,
            isAdmin: get().auth.isAdmin(),
            isEmailsEnabled: get().auth.hasEmailsPermission(),
            isTaskTrackersEnabled: get().auth.hasTaskTrackersPermission(),
            isNotesEnabled: get().auth.hasNotesPermission(),
            photoUrl: get().auth.getPhotoUrl(),
        })
    },
}))

export default useUser
