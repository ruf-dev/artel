import {Toggle} from "@vervstack/chures"

import cls from "@/components/S3ToggleFields/S3ToggleFields.module.css"

interface S3ToggleFieldsProps {
    useSsl: boolean
    pathStyle: boolean
    disabled: boolean
    onUseSslChange: (v: boolean) => void
    onPathStyleChange: (v: boolean) => void
}

export default function S3ToggleFields(
    {useSsl, pathStyle, disabled, onUseSslChange, onPathStyleChange}: S3ToggleFieldsProps,
) {
    return (
        <div className={cls.S3ToggleFieldsContainer}>
            <Toggle
                checked={useSsl}
                onChange={onUseSslChange}
                disabled={disabled}
                label="Use SSL"
                className={cls.ToggleRow}
            />
            <Toggle
                checked={pathStyle}
                onChange={onPathStyleChange}
                disabled={disabled}
                label="Path-style addressing"
                className={cls.ToggleRow}
            />
            <p className={cls.ToggleCaption}>
                Required for Garage and most self-hosted S3-compatible stores.
            </p>
        </div>
    )
}
