import {useState} from "react"
import {Button} from "@vervstack/chures"

import ModeBar from "@/pages/notes/components/ModeBar/ModeBar.tsx"
import SaveStatusIndicator, { SaveStatus } from "@/pages/notes/components/SaveStatusIndicator/SaveStatusIndicator.tsx"
import BreadcrumbPath from "@/pages/notes/components/BreadcrumbBar/components/BreadcrumbPath/BreadcrumbPath.tsx"
import PencilIcon from "@/pages/notes/components/icons/PencilIcon.tsx"
import CopyIcon from "@/pages/notes/components/icons/CopyIcon.tsx"
import CheckIcon from "@/pages/notes/components/icons/CheckIcon.tsx"
import cls from "@/pages/notes/components/BreadcrumbBar/BreadcrumbBar.module.css"


type Mode = 'edit' | 'preview' | 'read'

interface BreadcrumbBarProps {
    path: string | null
    mode: Mode
    onModeChange: (mode: Mode) => void
    saveStatus: SaveStatus
    saveError?: string
    onRename?: () => void
}

export default function BreadcrumbBar(
    { path, mode, onModeChange, saveStatus, saveError, onRename }: BreadcrumbBarProps
) {
    const [copied, setCopied] = useState(false)

    function handleCopyPath() {
        if (!path) {
            return
        }
        navigator.clipboard.writeText(path).then(() => {
            setCopied(true)
            setTimeout(() => setCopied(false), 1500)
        })
    }

    return (
        <div className={cls.BreadcrumbBarContainer}>
            <div className={cls.LeftSlot}>
                {path && <BreadcrumbPath path={path} />}
                {path && (
                    <Button
                        variant="ghost"
                        className={cls.CopyPathBtn}
                        onClick={handleCopyPath}
                        title={copied ? "Copied!" : "Copy path"}
                    >
                        {copied ? <CheckIcon /> : <CopyIcon />}
                    </Button>
                )}
                {path && onRename && (
                    <Button variant="ghost" className={cls.RenameBtn} onClick={onRename} title="Rename / Move">
                        <PencilIcon />
                    </Button>
                )}
            </div>
            <div className={cls.CenterSlot}>
                <SaveStatusIndicator status={saveStatus} errorMessage={saveError} />
            </div>
            <div className={cls.RightSlot}>
                <ModeBar active={mode} onModeChange={onModeChange} />
            </div>
        </div>
    )
}
