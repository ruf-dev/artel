import {Outlet, useNavigate} from 'react-router-dom'

import cls from '@/app/routing/HomeLayout.module.css'

import Topbar from '@/segments/Topbar/Topbar.tsx'
import useUser from "@/hooks/user/User.ts";
import {useEffect} from "react";
import {Path} from "@/app/routing/Router.tsx";

export default function HomeLayout() {
    const navigate = useNavigate()
    const {auth} = useUser()

    useEffect(() => {
        if (!auth.isAuthenticated()) {
            navigate(Path.InitPage)
        }
    }, [auth, navigate])

    return (
        <div className={cls.HomeLayoutContainer}>
            <main className={cls.Content}>
                <Outlet/>
            </main>
            <Topbar/>

        </div>
    )
}
