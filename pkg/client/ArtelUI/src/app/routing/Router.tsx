import cls from "@/app/routing/Router.module.css"

import {Navigate, Route, Routes} from "react-router-dom"

import InitPage from "@/pages/init/InitPage.tsx"
import HomePage from "@/pages/home/HomePage.tsx"
import ErrorPage from "@/pages/error/ErrorPage.tsx"

import useUser from "@/hooks/user/User.ts"

export enum Path {
    InitPage = "/init",
    HomePage = "/",
}

export default function Router() {
    const {auth, login, logout} = useUser()

    return (
        <div className={cls.Root}>
            <div className={cls.Content}>
                <Routes>
                    <Route
                        path={Path.InitPage}
                        element={<InitPage auth={auth} onLogin={login}/>}
                        errorElement={<ErrorPage/>}
                    />

                    <Route
                        path={Path.HomePage}
                        element={<HomePage auth={auth} onLogout={logout}/>}
                        errorElement={<ErrorPage/>}
                    />

                    <Route
                        path={"*"}
                        element={<Navigate to={Path.HomePage} replace/>}
                        errorElement={<ErrorPage/>}
                    />
                </Routes>
            </div>
        </div>
    )
}
