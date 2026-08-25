import {ReactNode} from "react"
import {Button} from "@vervstack/chures"
import {useNavigate} from "react-router-dom"

import {Path} from "@/app/routing/Router.tsx"
import BackIcon from "@/pages/workbench/components/WorkbenchToolbar/components/BackIcon/BackIcon.tsx"
import cls from "@/pages/workbench/components/WorkbenchTopBarShell/WorkbenchTopBarShell.module.css"

interface Props {
    vaultName: string
    statusBadge?: ReactNode
    actions?: ReactNode
}

// Shared top-bar chrome (back button, vault name, sticky positioning) reused by
// the Docker workbench's WorkbenchToolbar and Simple Chat's SimpleChatTopBar —
// each mode supplies its own statusBadge/actions content.
export default function WorkbenchTopBarShell(props: Props) {
    const navigate = useNavigate()

    return (
        <div className={cls.WorkbenchTopBarShellContainer}>
            <div className={cls.LeftSection}>
                <Button
                    variant="secondary"
                    className={cls.BackButton}
                    onClick={() => navigate(Path.HomePage)}
                    aria-label="Back to vaults"
                    title="Back to vaults"
                >
                    <BackIcon/>
                </Button>
                <span className={cls.VaultName}>{props.vaultName}</span>
                {props.statusBadge}
            </div>
            {props.actions}
        </div>
    )
}
