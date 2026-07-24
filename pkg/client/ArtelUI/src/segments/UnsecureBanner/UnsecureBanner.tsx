import {useAppConfig} from "@/app/hooks/AppConfig.ts"
import cls from "@/segments/UnsecureBanner/UnsecureBanner.module.css"

export default function UnsecureBanner() {
    const noAuthEnabled = useAppConfig((state) => state.noAuthEnabled)
    const credsEncrypted = useAppConfig((state) => state.credsEncrypted)

    if (!noAuthEnabled && credsEncrypted) return null

    return (
        <>
            {noAuthEnabled && (
                <div className={cls.UnsecureBannerContainer}>
                    running unsecure instance — authentication is disabled
                </div>
            )}
            {!credsEncrypted && (
                <div className={cls.UnsecureBannerContainer}>
                    running in insecure mode — all credentials are stored as plain text. generate an
                    encryption key and set it as ENVIRONMENT_CREDS_ENCRYPTION_KEY
                </div>
            )}
        </>
    )
}
