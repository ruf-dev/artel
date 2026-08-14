import {Link} from "react-router-dom"

import {Path} from "@/app/routing/Router.tsx"
import useUser from "@/hooks/user/User.ts"
import cls from "@/pages/docs/components/DocsEmptyState/DocsEmptyState.module.css"

// Friendly empty state for `/docs` when the admin hasn't configured a default public
// vault yet — distinct from DocsPage's "This vault could not be found." error, which is
// for a bad/unpublished slug reached via `/docs/:slug`.
export default function DocsEmptyState() {
    const {auth, isAdmin} = useUser()

    return (
        <div className={cls.DocsEmptyStateContainer}>
            <div className={cls.Message}>
                No public documentation has been configured yet.
            </div>
            {auth.isAuthenticated() && isAdmin && (
                <Link to={`${Path.Admin}?tab=settings`} className={cls.AdminLink}>Go to admin page</Link>
            )}
        </div>
    )
}
