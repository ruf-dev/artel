import cls from "@/components/VaultCard/VaultCardLinks.module.css"

export default function VaultCardLinks() {
    return (
        <div className={cls.VaultCardLinksContainer}>
            <a href="#" className={cls.Link} onClick={(e) => e.stopPropagation()}>
                LiveSync
            </a>
            <a href="#" className={cls.Link} onClick={(e) => e.stopPropagation()}>
                Connect to MCP
            </a>
        </div>
    )
}
