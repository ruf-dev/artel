import {useEffect} from "react"
import {useNavigate} from "react-router-dom"

import cls from "@/pages/init/InitPage.module.css"
import useUser from "@/hooks/user/User.ts"
import {Path} from "@/app/routing/Router.tsx"
import {useDialog} from "@/app/hooks/Dialog.ts"
import LoginContent from "@/pages/init/components/LoginContent/LoginContent.tsx"

export default function InitPage() {
    const navigate = useNavigate()
    const {auth, login} = useUser()
    const {OpenDialog, LockClosing} = useDialog()

    useEffect(function () {
        if (auth.isAuthenticated()) {
            navigate(Path.HomePage)
            return
        }
        LockClosing()
        OpenDialog(<LoginContent login={login} navigate={navigate}/>)
    }, [])

    return <div className={cls.Root}/>
}
