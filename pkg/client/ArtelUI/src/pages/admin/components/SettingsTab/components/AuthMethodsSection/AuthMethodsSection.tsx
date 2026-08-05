import {Toggle} from "@vervstack/chures"

import cls from "@/pages/admin/components/SettingsTab/SettingsTab.module.css"
import {useBakeError} from "@/app/hooks/useErrorToast"

interface AuthMethodsSectionProps {
    passwordAuthEnabled: boolean
    telegramAuthEnabled: boolean
    onPasswordAuthToggle: (enabled: boolean) => void
    onTelegramAuthToggle: (enabled: boolean) => void
}

export default function AuthMethodsSection({
    passwordAuthEnabled,
    telegramAuthEnabled,
    onPasswordAuthToggle,
    onTelegramAuthToggle,
}: AuthMethodsSectionProps) {
    const bakeError = useBakeError()

    function handlePasswordAuthToggle(enabled: boolean) {
        // Guard: can't disable both methods
        if (enabled === false && !telegramAuthEnabled) {
            bakeError("At least one login method must be enabled", new Error())
            return
        }
        onPasswordAuthToggle(enabled)
    }

    function handleTelegramAuthToggle(enabled: boolean) {
        // Guard: can't disable both methods
        if (enabled === false && !passwordAuthEnabled) {
            bakeError("At least one login method must be enabled", new Error())
            return
        }
        onTelegramAuthToggle(enabled)
    }

    return (
        <section className={cls.Section}>
            <h2 className={cls.SectionTitle}>Authentication Methods</h2>
            <div className={cls.SettingsGroup}>
                <Toggle
                    checked={passwordAuthEnabled}
                    onChange={handlePasswordAuthToggle}
                    label="Password login enabled"
                    className={cls.ToggleRow}
                />
                <Toggle
                    checked={telegramAuthEnabled}
                    onChange={handleTelegramAuthToggle}
                    label="Telegram login enabled"
                    className={cls.ToggleRow}
                />
            </div>
        </section>
    )
}
