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
    const {OpenDialog, LockClosing, UnlockClosing, CloseDialog} = useDialog()

    useEffect(function () {
        if (auth.isAuthenticated()) {
            navigate(Path.HomePage)
            return
        }
        LockClosing()
        OpenDialog(<LoginContent login={login} navigate={navigate}/>)

        // Router.tsx can navigate away from /init (e.g. setupCompleted === false → /setup)
        // before the user ever logs in. The login dialog lives in the global Dialog store,
        // not InitPage's own subtree, so it survives that navigation unless we explicitly
        // tear it down here — otherwise it's left open (and still lock-closed) on top of
        // whatever page comes next.
        return () => {
            UnlockClosing()
            CloseDialog()
        }
    }, [])

    return <div className={cls.Root}/>
}
