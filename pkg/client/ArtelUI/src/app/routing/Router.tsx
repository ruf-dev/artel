import {useEffect} from "react"

import cls from "@/app/routing/Router.module.css"

import {Navigate, Route, Routes} from "react-router-dom"

import InitPage from "@/pages/init/InitPage.tsx"
import HomePage from "@/pages/home/HomePage.tsx"
import EmailsPage from "@/pages/emails/EmailsPage.tsx"
import ErrorPage from "@/pages/error/ErrorPage.tsx"
import Dialog from "@/pages/segments/Dialog.tsx"
import McpAuthPage from "@/pages/mcp-auth/McpAuthPage.tsx"
import ClosedAlphaPage from "@/pages/closed-alpha/ClosedAlphaPage.tsx"
import AdminPage from "@/pages/admin/AdminPage.tsx"
import {AuthService} from "@/processes/Auth.ts"
import useUser from "@/hooks/user/User.ts"

// eslint-disable-next-line react-refresh/only-export-components
export enum Path {
    InitPage = "/init",
    HomePage = "/",
    EmailsPage = "/emails",
    McpAuth = "/authorize",
    ClosedAlpha = "/closed-alpha",
    Admin = "/admin",
}

export default function Router() {
    const {auth, setUserInfo} = useUser()

    useEffect(() => {
        if (!auth.isAuthenticated() || auth.userInfo) return
        const svc = new AuthService()
        svc.login(auth.session!)
        svc.FetchUserInfo().then(setUserInfo).catch(() => {})
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    return (
        <div className={cls.Root}>
            <div className={cls.Content}>
                <Routes>
                    <Route
                        path={Path.InitPage}
                        element={<InitPage/>}
                        errorElement={<ErrorPage/>}
                    />

                    <Route
                        path={Path.HomePage}
                        element={<HomePage/>}
                        errorElement={<ErrorPage/>}
                    />

                    <Route
                        path={Path.EmailsPage}
                        element={<EmailsPage/>}
                        errorElement={<ErrorPage/>}
                    />

                    <Route
                        path={Path.McpAuth}
                        element={<McpAuthPage/>}
                        errorElement={<ErrorPage/>}
                    />

                    <Route
                        path={Path.ClosedAlpha}
                        element={<ClosedAlphaPage/>}
                        errorElement={<ErrorPage/>}
                    />

                    <Route
                        path={Path.Admin}
                        element={<AdminPage/>}
                        errorElement={<ErrorPage/>}
                    />

                    <Route
                        path={"*"}
                        element={<Navigate to={Path.HomePage} replace/>}
                        errorElement={<ErrorPage/>}
                    />
                </Routes>

                <Dialog/>
            </div>
        </div>
    )
}
