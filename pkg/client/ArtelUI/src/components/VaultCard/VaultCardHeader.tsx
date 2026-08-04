import {Button} from "@vervstack/chures"
import {Tooltip} from "react-tooltip"
import {Link} from "react-router-dom"

import cls from "@/components/VaultCard/VaultCardHeader.module.css"
import {VaultItem} from "@/app/api/artel/vaults.pb.ts"
import {Path} from "@/app/routing/Router.tsx"
import useUser from "@/hooks/user/User.ts"

interface Props {
    vault: VaultItem
    onEdit?: (id: string) => void
    onWorkbench?: (id: string) => void
    isWorkbenchAvailable?: boolean
}

const GEAR_ICON_PATH = "M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33"
    + " 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1"
    + " 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6"
    + " 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0"
    + " 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06"
    + " .06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"

export default function VaultCardHeader({vault, onEdit, onWorkbench, isWorkbenchAvailable}: Props) {
    const {isAdmin} = useUser()
    const workbenchTooltipId = `workbench-tooltip-${vault.id ?? ""}`

    return (
        <>
            <div className={cls.VaultCardHeaderContainer}>
                <h3 className={cls.Name}>{vault.name}</h3>
                <div className={cls.Actions}>
                    {onWorkbench && (
                        <span
                            className={cls.WorkbenchButtonWrapper}
                            data-tooltip-id={isWorkbenchAvailable === false ? workbenchTooltipId : undefined}
                        >
                            <Button
                                className={cls.IconBtn}
                                disabled={isWorkbenchAvailable === false}
                                onClick={(e) => {
                                    e.stopPropagation()
                                    onWorkbench(vault.id ?? "")
                                }}
                                title="Open Workbench"
                                aria-label="Open Workbench"
                            >
                                <svg
                                    viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor"
                                    strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round"
                                >
                                    <rect x="3" y="4" width="18" height="16" rx="2"/>
                                    <path d="M7 9.5 10 12 7 14.5"/>
                                    <line x1="11" y1="15.5" x2="15" y2="15.5"/>
                                </svg>
                            </Button>
                        </span>
                    )}
                    {onEdit && (
                        <Button
                            className={cls.IconBtn}
                            onClick={(e) => {
                                e.stopPropagation()
                                onEdit(vault.id ?? "")
                            }}
                            title="Vault settings"
                            aria-label="Vault settings"
                        >
                            <svg
                                viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor"
                                strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round"
                            >
                                <circle cx="12" cy="12" r="3"/>
                                <path d={GEAR_ICON_PATH}/>
                            </svg>
                        </Button>
                    )}
                </div>
            </div>
            {isWorkbenchAvailable === false && (
                <Tooltip id={workbenchTooltipId} clickable noArrow>
                    {isAdmin
                        ? (
                            <>
                                Docker isn&apos;t configured.{" "}
                                <Link to={`${Path.Admin}?tab=docker_api&dialog=true`}>
                                    Set it up in Admin
                                </Link>
                            </>
                        )
                        : "Workbench requires Docker to be configured. Contact an administrator."}
                </Tooltip>
            )}
        </>
    )
}
