import {useEffect} from "react"
import {useNavigate} from "react-router-dom"

import cls from "@/pages/tract-canvas/TractCanvasListPage.module.css"

import {Path} from "@/app/routing/Router.tsx"
import {useTracts} from "@/app/hooks/Tracts.ts"
import useUser from "@/hooks/user/User.ts"

import HeroSegment from "@/pages/tract-canvas/segments/HeroSegment/HeroSegment.tsx"
import ContentSegment from "@/pages/tract-canvas/segments/ContentSegment/ContentSegment.tsx"

export default function TractCanvasListPage() {
    const navigate = useNavigate()
    const {auth} = useUser()
    const {fetch} = useTracts()

    useEffect(() => {
        if (!auth.isAuthenticated()) {
            navigate(Path.InitPage)
        }
    }, [auth, navigate])

    useEffect(() => {
        if (auth.isAuthenticated()) {
            void fetch()
        }
    }, [auth, fetch])

    return (
        <div className={cls.Root}>
            <HeroSegment/>
            <ContentSegment/>
        </div>
    )
}
