import {useEffect} from "react"

import cls from "@/app/routing/Router.module.css"

import {Navigate, Route, Routes, useNavigate} from "react-router-dom"

import InitPage from "@/pages/init/InitPage.tsx"
import HomePage from "@/pages/home/HomePage.tsx"
import McpKeysPage from "@/pages/mcp-keys/McpKeysPage.tsx"
import ErrorPage from "@/pages/error/ErrorPage.tsx"
import McpAuthPage from "@/pages/mcp-auth/McpAuthPage.tsx"
import ClosedAlphaPage from "@/pages/closed-alpha/ClosedAlphaPage.tsx"
import AdminPage from "@/pages/admin/AdminPage.tsx"
import TaskTrackersPage from "@/pages/task-trackers/TaskTrackersPage.tsx"
import JoinVaultPage from "@/pages/join/JoinVaultPage.tsx"
import NotesPage from "@/pages/notes/NotesPage.tsx"
import ConnectionsPage from "@/pages/connections/ConnectionsPage.tsx"
import ToolboxPage from "@/pages/toolbox/ToolboxPage.tsx"
import GoogleOAuthCallbackPage from "@/pages/connections/GoogleOAuthCallbackPage.tsx"
import HomeLayout from "@/app/routing/HomeLayout.tsx"
import Dialog from "@/pages/segments/Dialog.tsx"
import {authService} from "@/processes/Auth.ts"
import useUser from "@/hooks/user/User.ts"
import {useServerStatus} from "@/app/hooks/ServerStatus.ts"
import ServiceUnavailablePage from "@/pages/service-unavailable/ServiceUnavailablePage.tsx"
import { Tooltip } from "react-tooltip"
import { Toaster } from "@vervstack/chures"
// eslint-disable-next-line react-refresh/only-export-components
export enum Path {
    InitPage = "/init",
    HomePage = "/",
    McpKeysPage = "/mcp_keys",
    TaskTrackersPage = "/task-trackers",
    NotesPage = "/notes",
    ConnectionsPage = "/connections",
    ToolboxPage = "/toolbox",
    GoogleOAuthCallback = "/connections/google/callback",
    McpAuth = "/authorize",
    ClosedAlpha = "/closed-alpha",
    Admin = "/admin",
    JoinVault = "/join/:token",
}

export const REDIRECT_AFTER_LOGIN_KEY = "artel_post_login_redirect"

export default function Router() {
    const navigate = useNavigate()

    const {auth, setUserInfo} = useUser()
    const isServerAvailable = useServerStatus()

    useEffect(() => {
        if (!auth.isAuthenticated()) return

        authService.FetchUserInfo().then(setUserInfo).catch(() => {})
        if (location.pathname === Path.InitPage) {
            navigate(Path.HomePage)
        }
    }, [])

    if (!isServerAvailable) {
        return (
            <div className={cls.Root}>
                <div className={cls.Content}>
                    <ServiceUnavailablePage/>
                </div>
            </div>
        )
    }

    return (
        <div className={cls.Root}>
            <div className={cls.Content}>
                <Routes>
                    <Route element={<HomeLayout/>}>
                        <Route path={Path.HomePage} element={<HomePage/>} errorElement={<ErrorPage/>}/>
                        <Route path={Path.McpKeysPage} element={<McpKeysPage/>} errorElement={<ErrorPage/>}/>
                        <Route path={Path.Admin} element={<AdminPage/>} errorElement={<ErrorPage/>}/>
                        <Route path={Path.TaskTrackersPage} element={<TaskTrackersPage/>} errorElement={<ErrorPage/>}/>
                        <Route path={Path.NotesPage} element={<NotesPage/>} errorElement={<ErrorPage/>}/>
                        <Route path={Path.ConnectionsPage} element={<ConnectionsPage/>} errorElement={<ErrorPage/>}/>
                        <Route path={Path.GoogleOAuthCallback} element={<GoogleOAuthCallbackPage/>} errorElement={<ErrorPage/>}/>
                        <Route path={Path.ToolboxPage} element={<ToolboxPage/>} errorElement={<ErrorPage/>}/>
                        <Route path={"*"} element={<Navigate to={Path.HomePage} replace/>}/>
                    </Route>

                    <Route path={Path.InitPage} element={<InitPage/>} errorElement={<ErrorPage/>}/>
                    <Route path={Path.McpAuth} element={<McpAuthPage/>} errorElement={<ErrorPage/>}/>
                    <Route path={Path.ClosedAlpha} element={<ClosedAlphaPage/>} errorElement={<ErrorPage/>}/>
                    <Route path={Path.JoinVault} element={<JoinVaultPage/>} errorElement={<ErrorPage/>}/>
                </Routes>
            </div>
            <Dialog/>
            <Toaster/>
            <Tooltip
                id="root-tooltip"
                className={cls.Tooltip}
                classNameArrow={cls.TooltipArrow}
            />
        </div>
    )
}
