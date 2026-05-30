import {useState} from "react"
import cls from "@/components/VaultCard/VaultCardLinks.module.css"
import FastSetupDialog from "@/components/VaultCard/FastSetupDialog.tsx"

export default function VaultCardLinks() {
    const [dialogOpen, setDialogOpen] = useState(false)

    return (
        <div className={cls.VaultCardLinksContainer}>
            <a href="#" className={cls.Link} onClick={(e) => e.stopPropagation()}>
                LiveSync
            </a>
            <a href="#" className={cls.Link} onClick={(e) => e.stopPropagation()}>
                Connect to MCP
            </a>
            <button
                type="button"
                className={`${cls.Link} ${cls.LinkBtn}`}
                onClick={(e) => { e.stopPropagation(); setDialogOpen(true) }}
            >
                Fast setup
            </button>
            {dialogOpen && <FastSetupDialog onClose={() => setDialogOpen(false)}/>}
        </div>
    )
}
