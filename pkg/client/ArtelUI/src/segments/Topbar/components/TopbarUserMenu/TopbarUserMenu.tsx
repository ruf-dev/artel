import {useEffect, useRef, useState} from "react"
import {createPortal} from "react-dom"
import {useNavigate} from "react-router-dom"

import {useIsMobileNav} from "@/app/hooks/useIsMobileNav.ts"
import {Path} from "@/app/routing/Router.tsx"
import useUser from "@/hooks/user/User.ts"
import TopbarUserMenuPill from "@/segments/Topbar/components/TopbarUserMenuPill/TopbarUserMenuPill.tsx"
import TopbarUserMenuList from "@/segments/Topbar/components/TopbarUserMenuList/TopbarUserMenuList.tsx"
import cls from "@/segments/Topbar/components/TopbarUserMenu/TopbarUserMenu.module.css"

interface MenuRect {
    top?: number
    bottom?: number
    right: number
}

export default function TopbarUserMenu() {
    const [menuOpen, setMenuOpen] = useState(false)
    const [menuRect, setMenuRect] = useState<MenuRect | null>(null)
    const wrapRef = useRef<HTMLDivElement>(null)
    const navigate = useNavigate()
    const isMobile = useIsMobileNav()
    const {isAdmin, logout, photoUrl} = useUser()

    // Topbar has backdrop-filter, which makes it a containing block for position:fixed
    // descendants — the Backdrop's inset:0 would only cover the header strip instead of
    // the viewport. Portal to document.body and position via getBoundingClientRect
    // instead, same as TemplateInput.tsx.
    useEffect(() => {
        if (!menuOpen) return

        function reposition() {
            const el = wrapRef.current
            if (!el) return
            const box = el.getBoundingClientRect()
            // Mobile topbar is pinned to the bottom of the screen, so the menu has to
            // grow upward from the pill instead of downward off the viewport.
            if (isMobile) {
                setMenuRect({bottom: window.innerHeight - box.top + 8, right: window.innerWidth - box.right})
            } else {
                setMenuRect({top: box.bottom + 8, right: window.innerWidth - box.right})
            }
        }

        reposition()
        window.addEventListener("scroll", reposition, true)
        window.addEventListener("resize", reposition)
        return () => {
            window.removeEventListener("scroll", reposition, true)
            window.removeEventListener("resize", reposition)
        }
    }, [menuOpen, isMobile])

    function handleLogout() {
        setMenuOpen(false)
        logout()
        navigate(Path.InitPage)
    }

    function handleAdmin() {
        setMenuOpen(false)
        navigate(Path.Admin)
    }

    function handleApiKeys() {
        setMenuOpen(false)
        navigate(Path.McpKeysPage)
    }

    function handleSkills() {
        setMenuOpen(false)
        navigate(Path.SkillsPage)
    }

    function handleDocs() {
        setMenuOpen(false)
        navigate(Path.DocsPageDefault)
    }

    return (
        <div className={cls.UserWrap} ref={wrapRef}>
            <TopbarUserMenuPill menuOpen={menuOpen} photoUrl={photoUrl} onClick={e => {
                e.stopPropagation();
                setMenuOpen(v => !v)
            }}/>
            {menuOpen && menuRect && createPortal(
                <>
                    <div className={cls.Backdrop} onClick={() => setMenuOpen(false)}/>
                    <TopbarUserMenuList
                        rect={menuRect}
                        isAdmin={isAdmin}
                        onAdmin={handleAdmin}
                        onApiKeys={handleApiKeys}
                        onSkills={handleSkills}
                        onDocs={handleDocs}
                        onLogout={handleLogout}
                    />
                </>,
                document.body,
            )}
        </div>
    )
}
