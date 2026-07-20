import {useEffect} from "react"
import {Navigate, useNavigate, useRoutes, type RouteObject} from "react-router-dom"
import { Tooltip } from "react-tooltip"
import { Toaster } from "@vervstack/chures"

import cls from "@/app/routing/Router.module.css"
import InitPage from "@/pages/init/InitPage.tsx"
import HomePage from "@/pages/home/HomePage.tsx"
import McpKeysPage from "@/pages/mcp-keys/McpKeysPage.tsx"
import ErrorPage from "@/pages/error/ErrorPage.tsx"
import McpAuthPage from "@/pages/mcp-auth/McpAuthPage.tsx"
import ClosedAlphaPage from "@/pages/closed-alpha/ClosedAlphaPage.tsx"
import AdminPage from "@/pages/admin/AdminPage.tsx"
import TaskTrackersPage from "@/pages/task-trackers/TaskTrackersPage.tsx"
import RoadmapPage from "@/pages/roadmap/RoadmapPage.tsx"
import JoinVaultPage from "@/pages/join/JoinVaultPage.tsx"
import NotesPage from "@/pages/notes/NotesPage.tsx"
import ConnectionsPage from "@/pages/connections/ConnectionsPage.tsx"
import ToolboxPage from "@/pages/toolbox/ToolboxPage.tsx"
import TractCanvasListPage from "@/pages/tract-canvas/TractCanvasListPage.tsx"
import TractCanvasBuilderPage from "@/pages/tract-canvas/TractCanvasBuilderPage.tsx"
import TractTemplatesListPage from "@/pages/tract-templates/TractTemplatesListPage.tsx"
import GoogleOAuthCallbackPage from "@/pages/connections/GoogleOAuthCallbackPage.tsx"
import HomeLayout from "@/app/routing/HomeLayout.tsx"
import Dialog from "@/pages/segments/Dialog.tsx"
import {authService} from "@/processes/Auth.ts"
import useUser from "@/hooks/user/User.ts"
import {useServerStatus} from "@/app/hooks/ServerStatus.ts"
import ServiceUnavailablePage from "@/pages/service-unavailable/ServiceUnavailablePage.tsx"

// eslint-disable-next-line react-refresh/only-export-components
export enum Path {
    InitPage = "/init",
    HomePage = "/",
    McpKeysPage = "/mcp_keys",
    TaskTrackersPage = "/task-trackers",
    RoadmapPage = "/task-trackers/:trackerId/roadmap/:cardId",
    NotesPage = "/notes",
    NotesPageVault = "/notes/:vaultId",
    NotesPageNote = "/notes/:vaultId/*",
    ConnectionsPage = "/connections",
    ToolboxPage = "/toolbox",
    TractsPage = "/tracts",
    TractEditorPage = "/tracts/:tractUuid",
    TractTemplatesPage = "/tract-templates",
    GoogleOAuthCallback = "/connections/google/callback",
    McpAuth = "/authorize",
    ClosedAlpha = "/closed-alpha",
    Admin = "/admin",
    JoinVault = "/join/:token",
}

export const REDIRECT_AFTER_LOGIN_KEY = "artel_post_login_redirect"

const routes: RouteObject[] = [
    {
        element: <HomeLayout/>,
        children: [
            {path: Path.HomePage, element: <HomePage/>, errorElement: <ErrorPage/>},
            {path: Path.McpKeysPage, element: <McpKeysPage/>, errorElement: <ErrorPage/>},
            {path: Path.Admin, element: <AdminPage/>, errorElement: <ErrorPage/>},
            {path: Path.TaskTrackersPage, element: <TaskTrackersPage/>, errorElement: <ErrorPage/>},
            {path: Path.RoadmapPage, element: <RoadmapPage/>, errorElement: <ErrorPage/>},
            {path: Path.NotesPage, element: <NotesPage/>, errorElement: <ErrorPage/>},
            {path: Path.NotesPageVault, element: <NotesPage/>, errorElement: <ErrorPage/>},
            {path: Path.NotesPageNote, element: <NotesPage/>, errorElement: <ErrorPage/>},
            {path: Path.ConnectionsPage, element: <ConnectionsPage/>, errorElement: <ErrorPage/>},
            {path: Path.GoogleOAuthCallback, element: <GoogleOAuthCallbackPage/>, errorElement: <ErrorPage/>},
            {path: Path.ToolboxPage, element: <ToolboxPage/>, errorElement: <ErrorPage/>},
            {path: Path.TractsPage, element: <TractCanvasListPage/>, errorElement: <ErrorPage/>},
            {path: Path.TractEditorPage, element: <TractCanvasBuilderPage/>, errorElement: <ErrorPage/>},
            {path: Path.TractTemplatesPage, element: <TractTemplatesListPage/>, errorElement: <ErrorPage/>},
            {path: "*", element: <Navigate to={Path.HomePage} replace/>},
        ],
    },
    {path: Path.InitPage, element: <InitPage/>, errorElement: <ErrorPage/>},
    {path: Path.McpAuth, element: <McpAuthPage/>, errorElement: <ErrorPage/>},
    {path: Path.ClosedAlpha, element: <ClosedAlphaPage/>, errorElement: <ErrorPage/>},
    {path: Path.JoinVault, element: <JoinVaultPage/>, errorElement: <ErrorPage/>},
]

export default function Router() {
    const navigate = useNavigate()

    const {auth, setUserInfo} = useUser()
    const isServerAvailable = useServerStatus()
    const routeElement = useRoutes(routes)

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
                {routeElement}
            </div>
            <Dialog/>
            <Toaster/>
            <Tooltip
                id="root-tooltip"
                className={cls.Tooltip}
                noArrow
            />
        </div>
    )
}
