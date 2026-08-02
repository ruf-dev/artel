// eslint-disable-next-line max-len -- import path too long to wrap under 120 chars
import cls from "@/pages/notes/components/NotesSidebar/components/SidebarTopBar/components/PublicBadge/PublicBadge.module.css"

export default function PublicBadge() {
    return (
        <span className={cls.PublicBadgeContainer}>Public</span>
    )
}
