import cls from "@/app/routing/Router.module.css"

import {Navigate, Route, Routes} from "react-router-dom"

import InitPage from "@/pages/init/InitPage.tsx"
import HomePage from "@/pages/home/HomePage.tsx"
import ErrorPage from "@/pages/error/ErrorPage.tsx"
import Dialog from "@/pages/segments/Dialog.tsx"
import McpAuthPage from "@/pages/mcp-auth/McpAuthPage.tsx"

// eslint-disable-next-line react-refresh/only-export-components
export enum Path {
    InitPage = "/init",
    HomePage = "/",
    McpAuth = "/authorize",
}

export default function Router() {
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
                        path={Path.McpAuth}
                        element={<McpAuthPage/>}
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
