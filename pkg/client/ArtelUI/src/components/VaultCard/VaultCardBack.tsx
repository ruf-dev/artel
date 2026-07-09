import React from "react"
import {Button, InfoDialog} from "@vervstack/chures"

import cls from "@/components/VaultCard/VaultCardBack.module.css"
import FastSetupDialog from "@/dialogs/FastSetupDialog/FastSetupDialog.tsx";
import {useDialog} from "@/app/hooks/Dialog.ts";
import {cn} from "@/app/utils/cn";

export default function VaultCardBack() {

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
