import React, {useState} from "react"
import {Button, InfoDialog} from "@vervstack/chures"
import cls from "@/widgets/VaultCard/VaultCard.module.css"

import {VaultItem} from "@/app/api/artel/vaults.pb.ts"

import VaultCardHeader from "@/components/VaultCard/VaultCardHeader.tsx";
import VaultCardStatus from "@/components/VaultCard/VaultCardStatus.tsx";
import VaultCardConnBar from "@/components/VaultCard/VaultCardConnBar.tsx";
import FastSetupDialog from "@/dialogs/FastSetupDialog/FastSetupDialog.tsx";
import {useDialog} from "@/app/hooks/Dialog.ts";
import {cn} from "@/app/utils/cn";

interface Props {
    vault: VaultItem
    onEdit?: (id: string) => void
}

export default function VaultCard({vault, onEdit}: Props) {
    const [isFlipped, setIsFlipped] = useState(false)

    return (
        <article className={cls.VaultCardContainer} onClick={() => setIsFlipped(!isFlipped)}>
            <div className={cn(cls.VaultCardContentWrapper, isFlipped && cls.Flipped)}>
                <VaultCardFront vault={vault} onEdit={onEdit}/>
                <VaultCardBack/>
            </div>
        </article>
    )
}

function VaultCardFront({vault, onEdit}: Props) {
    return (
        <div className={cls.VaultCardFrontContainer}>
            <VaultCardHeader vault={vault} onEdit={onEdit}/>
            <VaultCardStatus/>
            <VaultCardConnBar dbUrl={vault.dbUrl ?? ""} passphrase={vault.livesyncPassphrase ?? ""}/>
        </div>
    )
}


function VaultCardBack() {

    const {OpenDialog, CloseDialog} = useDialog()

    function openSetupDialog(e: React.MouseEvent<HTMLButtonElement>) {
        e.stopPropagation()
        OpenDialog(<FastSetupDialog onClose={CloseDialog}/>)
    }

    function openMcpSetupDialog(e: React.MouseEvent<HTMLButtonElement>) {
        e.stopPropagation()
        OpenDialog(<InfoDialog
            title={'Not implemented yet'}
            message={'Later there will be instructions on how to setup different AI agents'}
            onClose={CloseDialog}
        />)

    }

    return (
        <div className={cls.VaultCardBackContainer}>
            <div className={cls.VaultCardLinksContainer}>
                <Button
                    className={cls.Link}
                    onClick={openMcpSetupDialog}>
                    Connect to MCP
                </Button>
                <Button
                    className={cn(cls.Link, cls.LinkBtn)}
                    onClick={openSetupDialog}>
                    Fast setup
                </Button>

            </div>
        </div>
    )
}
