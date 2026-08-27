import {useState, useEffect} from "react"
import {Button, ConfirmDialog} from "@vervstack/chures"

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
import {useDialog} from "@/app/hooks/Dialog.ts"

// eslint-disable-next-line max-lines-per-function
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
    const [systemPrompt, setSystemPrompt] = useState("")
    const {OpenDialog, CloseDialog} = useDialog()

    useEffect(() => {
        AdminSystemSettingsAPI.GetSettings({}, auth.getInitReq())
            .then(res => {
                setPasswordAuthEnabled(res.passwordAuthEnabled ?? false)
                setTelegramAuthEnabled(res.telegramAuthEnabled ?? false)
                setRegistrationMode(res.registrationMode ?? RegistrationMode.ADMIN_ONLY)
                setDefaultDocsVaultId(res.defaultDocsVaultId ?? "")
                setDocsSource(res.defaultDocsSource ?? DocsSource.VAULT)
                setSystemPrompt(res.systemPrompt ?? "")
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

    function saveSystemPrompt() {
        OpenDialog(
            <ConfirmDialog
                title="Change system prompt"
                message="You are about to change system wide promt. it will affect ALL users. Are you sure?"
                confirmLabel="Save prompt"
                onClose={CloseDialog}
                onConfirm={() => AdminSystemSettingsAPI.UpdateSystemPrompt({prompt: systemPrompt}, auth.getInitReq())
                    .then(() => undefined)
                    .catch(err => bakeError("Failed to update system prompt", err))}
            />,
        )
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
            <section className={cls.Section}>
                <h2 className={cls.SectionTitle}>System prompt</h2>
                <div className={cls.SettingsGroup}>
                    <textarea className={cls.PromptInput} value={systemPrompt}
                        onChange={e => setSystemPrompt(e.target.value)} rows={8} />
                    <Button className={cls.PromptButton} onClick={saveSystemPrompt}>Save system prompt</Button>
                </div>
            </section>
        </div>
    )
}
