import {useEffect} from 'react'
import {Outlet, useLocation, useNavigate} from 'react-router-dom'

import cls from '@/app/routing/HomeLayout.module.css'
import Topbar from '@/segments/Topbar/Topbar.tsx'
import useUser from '@/hooks/user/User.ts'
import {Path, REDIRECT_AFTER_LOGIN_KEY} from '@/app/routing/Router.tsx'

export default function HomeLayout() {
	const {auth} = useUser()
	const location = useLocation()
	const navigate = useNavigate()

	useEffect(() => {
		if (auth.isAuthenticated()) return
		localStorage.setItem(REDIRECT_AFTER_LOGIN_KEY, location.pathname)
		navigate(Path.InitPage)
	}, [location.pathname])

	return (
		<div className={cls.HomeLayoutContainer}>
			<main className={cls.Content}>
				<Outlet/>
			</main>
			<Topbar/>
		</div>
	)
}
