import {Toggle} from "@vervstack/chures"

import cls from
"@/dialogs/ManageVaultDialog/widgets/VaultSettingsSection/components/BinaryStorageToggle/BinaryStorageToggle.module.css"

interface Props {
    checked: boolean
    disabled: boolean
    onChange: (v: boolean) => void
}

export default function BinaryStorageToggle({checked, disabled, onChange}: Props) {
    return (
        <div className={cls.BinaryStorageToggleContainer}>
            <Toggle
                checked={checked}
                onChange={onChange}
                disabled={disabled}
                label="Use for binary data"
            />
            <p className={cls.Caption}>
                When on, binary files (images, PDFs, etc.) are stored in CouchDB. When off,
                they&apos;re stored in the vault&apos;s linked S3 bucket.
            </p>
        </div>
    )
}
