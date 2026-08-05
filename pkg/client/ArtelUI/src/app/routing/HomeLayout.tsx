import {Outlet} from 'react-router-dom'

import cls from '@/app/routing/HomeLayout.module.css'
import Topbar from '@/segments/Topbar/Topbar.tsx'

export default function HomeLayout() {
    return (
        <div className={cls.HomeLayoutContainer}>
            <main className={cls.Content}>
                <Outlet/>
            </main>
            <Topbar/>

        </div>
    )
}
