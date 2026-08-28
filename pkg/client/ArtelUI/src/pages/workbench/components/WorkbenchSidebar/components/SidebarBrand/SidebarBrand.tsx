import {useNavigate} from "react-router-dom"
import {Button} from "@vervstack/chures"

import {Path} from "@/app/routing/Router.tsx"
import CloseIcon from "@/icons/common/CloseIcon.tsx"
import cls from "@/pages/workbench/components/WorkbenchSidebar/components/SidebarBrand/SidebarBrand.module.css"
import SidebarBrandMark from "@/pages/workbench/components/WorkbenchSidebar/components/icons/SidebarBrandMark.tsx"

interface Props {
    showClose?: boolean
}

// The sidebar's top row: the artel wordmark (navigates Home) plus, in api mode
// only, an "exit to home" control. The close control rides here temporarily —
// Stage 3 relocates it.
export default function SidebarBrand({showClose}: Props) {
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
            {showClose && (
                <Button
                    variant="secondary"
                    className={cls.CloseChatButton}
                    onClick={() => navigate(Path.HomePage)}
                    aria-label="Close chat"
                    title="Exit to home"
                >
                    <CloseIcon className={cls.CloseChatIcon}/>
                </Button>
            )}
        </div>
    )
}
