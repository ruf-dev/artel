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
import DocsPage from "@/pages/docs/DocsPage.tsx"
import DocsDefaultRedirect from "@/pages/docs/DocsDefaultRedirect.tsx"
import ConnectionsPage from "@/pages/connections/ConnectionsPage.tsx"
import ToolboxPage from "@/pages/toolbox/ToolboxPage.tsx"
import TractCanvasListPage from "@/pages/tract-canvas/TractCanvasListPage.tsx"
import TractCanvasBuilderPage from "@/pages/tract-canvas/TractCanvasBuilderPage.tsx"
import TractTemplatesListPage from "@/pages/tract-templates/TractTemplatesListPage.tsx"
import GoogleOAuthCallbackPage from "@/pages/connections/GoogleOAuthCallbackPage.tsx"
import WorkbenchPage from "@/pages/workbench/WorkbenchPage.tsx"
import HomeLayout from "@/app/routing/HomeLayout.tsx"
import Dialog from "@/pages/segments/Dialog.tsx"
import UnsecureBanner from "@/segments/UnsecureBanner/UnsecureBanner.tsx"
import {authService} from "@/processes/Auth.ts"
import useUser from "@/hooks/user/User.ts"
import {useServerStatus} from "@/app/hooks/ServerStatus.ts"
import {useAppConfig} from "@/app/hooks/AppConfig.ts"
import ServiceUnavailablePage from "@/pages/service-unavailable/ServiceUnavailablePage.tsx"
import SetupWizardPage from "@/pages/setup-wizard/SetupWizardPage.tsx"

// eslint-disable-next-line react-refresh/only-export-components
export enum Path {
    InitPage = "/init",
    SetupWizard = "/setup",
    HomePage = "/",
    McpKeysPage = "/mcp_keys",
    TaskTrackersPage = "/task-trackers",
    RoadmapPage = "/task-trackers/:trackerId/roadmap/:cardId",
    NotesPage = "/notes",
    NotesPageVault = "/notes/:vaultId",
    NotesPageNote = "/notes/:vaultId/*",
    WorkbenchPage = "/workbench/:vaultId",
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
    DocsPage = "/docs/:slug",
    DocsPageDefault = "/docs",
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
            {path: Path.WorkbenchPage, element: <WorkbenchPage/>, errorElement: <ErrorPage/>},
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
    {path: Path.SetupWizard, element: <SetupWizardPage/>, errorElement: <ErrorPage/>},
    {path: Path.McpAuth, element: <McpAuthPage/>, errorElement: <ErrorPage/>},
    {path: Path.ClosedAlpha, element: <ClosedAlphaPage/>, errorElement: <ErrorPage/>},
    {path: Path.JoinVault, element: <JoinVaultPage/>, errorElement: <ErrorPage/>},
    {path: Path.DocsPage, element: <DocsPage/>, errorElement: <ErrorPage/>},
    {
        path: Path.DocsPageDefault,
        element: <DocsDefaultRedirect/>,
        errorElement: <ErrorPage/>,
    },
]

export default function Router() {
    const navigate = useNavigate()

    const {auth, setUserInfo} = useUser()
    const isServerAvailable = useServerStatus()
    const setupCompleted = useAppConfig((s) => s.setupCompleted)
    const routeElement = useRoutes(routes)

    // setupCompleted starts at its safe default (true) and only reflects the real value once
    // the first GetConfig response lands — this must re-run when that happens, not just once at
    // mount, or a fresh instance never actually gets routed to /setup.
    useEffect(() => {
        if (!setupCompleted && location.pathname !== Path.SetupWizard) {
            navigate(Path.SetupWizard)
        }
    }, [setupCompleted])

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
            <UnsecureBanner/>
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
