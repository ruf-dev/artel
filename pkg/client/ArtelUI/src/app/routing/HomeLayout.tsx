import {Outlet} from 'react-router-dom'
import {Tooltip} from 'react-tooltip'

import cls from '@/app/routing/HomeLayout.module.css'

import Topbar from '@/segments/Topbar/Topbar.tsx'
import Dialog from '@/pages/segments/Dialog.tsx'

export default function HomeLayout() {
    return (
        <div className={cls.HomeLayoutContainer}>
            <Topbar/>
            <main className={cls.Content}>
                <Outlet/>
            </main>
            <Dialog/>
            <Tooltip id="root-tooltip"/>
        </div>
    )
}
