import {useState, useEffect} from "react"

import {
    AdminSystemSettingsAPI,
} from "@/app/api/artel/admin_system_settings.pb.ts"
import {DocsSource} from "@/app/api/artel/public_docs.pb.ts"
import {RegistrationMode} from "@/app/api/artel/setup_wizard.pb.ts"
import {useBakeError} from "@/app/hooks/useErrorToast"
import useUser from "@/hooks/user/User.ts"
import cls from "@/pages/admin/components/SettingsTab/SettingsTab.module.css"
import AuthMethodsSection from "@/pages/admin/components/SettingsTab/components/AuthMethodsSection/AuthMethodsSection.tsx" // eslint-disable-line max-len
import RegistrationModeSection from "@/pages/admin/components/SettingsTab/components/RegistrationModeSection/RegistrationModeSection.tsx" // eslint-disable-line max-len
import DefaultDocsVaultSection from "@/pages/admin/components/SettingsTab/components/DefaultDocsVaultSection/DefaultDocsVaultSection.tsx" // eslint-disable-line max-len
import {createDocsSettingsHandlers} from "@/pages/admin/components/SettingsTab/processes/docsSettingsHandlers.ts"

export default function SettingsTab() {
    const {auth} = useUser()
    const bakeError = useBakeError()
    const [loading, setLoading] = useState(true)
    const [passwordAuthEnabled, setPasswordAuthEnabled] = useState(false)
    const [telegramAuthEnabled, setTelegramAuthEnabled] = useState(false)
    const [registrationMode, setRegistrationMode] = useState<RegistrationMode>(
        RegistrationMode.ADMIN_ONLY,
    )
    const [defaultDocsVaultId, setDefaultDocsVaultId] = useState("")
    const [docsSource, setDocsSource] = useState<DocsSource>(DocsSource.VAULT)

    useEffect(() => {
        AdminSystemSettingsAPI.GetSettings({}, auth.getInitReq())
            .then(res => {
                setPasswordAuthEnabled(res.passwordAuthEnabled ?? false)
                setTelegramAuthEnabled(res.telegramAuthEnabled ?? false)
                setRegistrationMode(res.registrationMode ?? RegistrationMode.ADMIN_ONLY)
                setDefaultDocsVaultId(res.defaultDocsVaultId ?? "")
                setDocsSource(res.defaultDocsSource ?? DocsSource.VAULT)
            })
            .catch(err => bakeError("Failed to load settings", err))
            .finally(() => setLoading(false))
    }, [auth, bakeError])

    function handlePasswordAuthToggle(enabled: boolean) {
        const oldValue = passwordAuthEnabled
        setPasswordAuthEnabled(enabled)

        AdminSystemSettingsAPI.UpdateAuthMethods(
            {passwordEnabled: enabled, telegramEnabled: telegramAuthEnabled},
            auth.getInitReq(),
        )
            .catch(err => {
                bakeError("Failed to update auth methods", err)
                setPasswordAuthEnabled(oldValue)
            })
    }

    function handleTelegramAuthToggle(enabled: boolean) {
        const oldValue = telegramAuthEnabled
        setTelegramAuthEnabled(enabled)

        AdminSystemSettingsAPI.UpdateAuthMethods(
            {passwordEnabled: passwordAuthEnabled, telegramEnabled: enabled},
            auth.getInitReq(),
        )
            .catch(err => {
                bakeError("Failed to update auth methods", err)
                setTelegramAuthEnabled(oldValue)
            })
    }

    function handleRegistrationModeChange(mode: RegistrationMode) {
        if (registrationMode === mode) return

        const oldMode = registrationMode
        setRegistrationMode(mode)

        AdminSystemSettingsAPI.UpdateRegistrationMode({mode}, auth.getInitReq())
            .catch(err => {
                bakeError("Failed to update registration mode", err)
                setRegistrationMode(oldMode)
            })
    }

    const {handleDefaultDocsVaultChange, handleDocsSourceChange} = createDocsSettingsHandlers({
        auth,
        defaultDocsVaultId,
        docsSource,
        setDefaultDocsVaultId,
        setDocsSource,
        bakeError,
    })

    if (loading) {
        return <div className={cls.SettingsTabContainer}>Loading settings…</div>
    }

    return (
        <div className={cls.SettingsTabContainer}>
            <AuthMethodsSection
                passwordAuthEnabled={passwordAuthEnabled}
                telegramAuthEnabled={telegramAuthEnabled}
                onPasswordAuthToggle={handlePasswordAuthToggle}
                onTelegramAuthToggle={handleTelegramAuthToggle}
            />
            <RegistrationModeSection
                registrationMode={registrationMode}
                onRegistrationModeChange={handleRegistrationModeChange}
            />
            <DefaultDocsVaultSection
                defaultDocsVaultId={defaultDocsVaultId}
                onDefaultDocsVaultChange={handleDefaultDocsVaultChange}
                docsSource={docsSource}
                onDocsSourceChange={handleDocsSourceChange}
            />
        </div>
    )
}
