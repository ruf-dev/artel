import {Toggle} from "@vervstack/chures"

import cls from
    "@/dialogs/ManageVaultDialog/widgets/VaultSettingsSection/components/PublishToggle/PublishToggle.module.css"

interface Props {
    checked: boolean
    disabled: boolean
    onChange: (v: boolean) => void
}

export default function PublishToggle({checked, disabled, onChange}: Props) {
    return (
        <div className={cls.PublishToggleContainer}>
            <Toggle
                checked={checked}
                onChange={onChange}
                disabled={disabled}
                label="Make public"
            />
            <p className={cls.Caption}>
                When on, anyone with the public link can view this vault&apos;s contents without signing in.
            </p>
        </div>
    )
}
