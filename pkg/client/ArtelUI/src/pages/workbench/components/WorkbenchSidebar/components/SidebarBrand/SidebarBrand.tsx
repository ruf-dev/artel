import {useNavigate} from "react-router-dom"
import {Button} from "@vervstack/chures"

import {Path} from "@/app/routing/Router.tsx"
import cls from "@/pages/workbench/components/WorkbenchSidebar/components/SidebarBrand/SidebarBrand.module.css"
import SidebarBrandMark from "@/pages/workbench/components/WorkbenchSidebar/components/icons/SidebarBrandMark.tsx"

// The sidebar's top row: the artel wordmark, which navigates Home.
export default function SidebarBrand() {
    const navigate = useNavigate()

    return (
        <div className={cls.SidebarBrandContainer}>
            <Button
                variant="unstyled"
                className={cls.BrandButton}
                onClick={() => navigate(Path.HomePage)}
                aria-label="Go to home"
            >
                <SidebarBrandMark className={cls.BrandMark}/>
                <span className={cls.Wordmark}>artel</span>
            </Button>
        </div>
    )
}
