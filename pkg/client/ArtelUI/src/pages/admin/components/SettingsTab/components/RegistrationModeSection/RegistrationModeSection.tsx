import {Button} from "@vervstack/chures"

import cls from "@/pages/admin/components/SettingsTab/SettingsTab.module.css"
import {RegistrationMode} from "@/app/api/artel/setup_wizard.pb.ts"
import {cn} from "@/app/utils/cn.ts"

interface RegistrationModeSectionProps {
    registrationMode: RegistrationMode
    onRegistrationModeChange: (mode: RegistrationMode) => void
}

export default function RegistrationModeSection({
    registrationMode,
    onRegistrationModeChange,
}: RegistrationModeSectionProps) {
    const isAdminOnly = registrationMode === RegistrationMode.ADMIN_ONLY
    const isSelfRegister = registrationMode === RegistrationMode.SELF_REGISTER

    return (
        <section className={cls.Section}>
            <h2 className={cls.SectionTitle}>Registration Mode</h2>
            <div className={cls.ModeButtonGroup}>
                <Button
                    variant={isAdminOnly ? "primary" : "secondary"}
                    onClick={() => onRegistrationModeChange(RegistrationMode.ADMIN_ONLY)}
                    className={cn(cls.ModeButton, isAdminOnly && cls.ModeButtonActive)}
                >
                    Admin Only
                </Button>
                <Button
                    variant={isSelfRegister ? "primary" : "secondary"}
                    onClick={() => onRegistrationModeChange(RegistrationMode.SELF_REGISTER)}
                    className={cn(cls.ModeButton, isSelfRegister && cls.ModeButtonActive)}
                >
                    Self Register
                </Button>
            </div>
        </section>
    )
}
