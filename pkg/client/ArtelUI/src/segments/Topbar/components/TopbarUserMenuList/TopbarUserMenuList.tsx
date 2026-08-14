import {cn} from "@/app/utils/cn.ts"
import TopbarUserMenuItem from "@/segments/Topbar/components/TopbarUserMenuItem/TopbarUserMenuItem.tsx"
import cls from "@/segments/Topbar/components/TopbarUserMenuList/TopbarUserMenuList.module.css"

interface MenuRect {
    top?: number
    bottom?: number
    right: number
}

interface TopbarUserMenuListProps {
    isAdmin: boolean
    onAdmin: () => void
    onApiKeys: () => void
    onDocs: () => void
    onLogout: () => void
    rect: MenuRect
}

export default function TopbarUserMenuList(
    {isAdmin, onAdmin, onApiKeys, onDocs, onLogout, rect}: TopbarUserMenuListProps,
) {
    return (
        <div className={cn(cls.Menu, rect.bottom !== undefined && cls.MenuUpward)} role="menu" style={rect}>
            {isAdmin && (
                <TopbarUserMenuItem
                    label="Admin Panel"
                    onClick={onAdmin}
                    icon={
                        <svg viewBox="0 0 24 24" width="15" height="15" fill="none" stroke="currentColor"
                             strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
                            <circle cx="12" cy="12" r="3"/>
                            <path d="M19.07 4.93a10 10 0 0 1 0 14.14M4.93 4.93a10 10 0 0 0 0 14.14"/>
                        </svg>
                    }
                />
            )}
            <TopbarUserMenuItem
                label="API Keys"
                onClick={onApiKeys}
                icon={
                    <svg viewBox="0 0 24 24" width="15" height="15" fill="none" stroke="currentColor"
                         strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
                        {/* eslint-disable-next-line max-len -- SVG path data can't be wrapped */}
                        <path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/>
                    </svg>
                }
            />
            <TopbarUserMenuItem
                label="Docs"
                onClick={onDocs}
                icon={
                    <svg viewBox="0 0 24 24" width="15" height="15" fill="none" stroke="currentColor"
                         strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
                        <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"/>
                        <polyline points="13 2 13 9 20 9"/>
                        <line x1="9" y1="15" x2="15" y2="15"/>
                        <line x1="9" y1="19" x2="15" y2="19"/>
                    </svg>
                }
            />
            <TopbarUserMenuItem
                label="Log out"
                onClick={onLogout}
                danger
                icon={
                    <svg viewBox="0 0 24 24" width="15" height="15" fill="none" stroke="currentColor"
                         strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
                        <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
                        <polyline points="16 17 21 12 16 7"/>
                        <line x1="21" y1="12" x2="9" y2="12"/>
                    </svg>
                }
            />
        </div>
    )
}
