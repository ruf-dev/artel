import {useState} from "react"
import {Button, ChevronDownIcon} from "@vervstack/chures"

import cls from "@/dialogs/ManageVaultDialog/widgets/ExpertSettingsSection/ExpertSettingsSection.module.css"
import {VaultItem} from "@/app/api/artel/vaults.pb.ts"
import ExpertSettingsDrawer from "@/dialogs/ManageVaultDialog/components/ExpertSettingsDrawer/ExpertSettingsDrawer.tsx"

interface Props {
    vault: VaultItem
    onChanged: (patch: Partial<VaultItem>) => void
}

export default function ExpertSettingsSection({vault, onChanged}: Props) {
    const [open, setOpen] = useState(false)

    return (
        <section className={cls.ExpertSettingsSectionContainer}>
            <Button variant="ghost" className={cls.TriggerBtn} onClick={() => setOpen(true)}>
                <span>Expert settings</span>
                <ChevronDownIcon className={cls.Chevron}/>
            </Button>
            <ExpertSettingsDrawer
                open={open}
                vault={vault}
                onChanged={onChanged}
                onClose={() => setOpen(false)}
            />
        </section>
    )
}
