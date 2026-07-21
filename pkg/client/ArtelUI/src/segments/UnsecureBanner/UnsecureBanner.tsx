import {useAppConfig} from "@/app/hooks/AppConfig.ts"
import cls from "@/segments/UnsecureBanner/UnsecureBanner.module.css"

export default function UnsecureBanner() {
    const noAuthEnabled = useAppConfig((state) => state.noAuthEnabled)

    if (!noAuthEnabled) return null

    return (
        <div className={cls.UnsecureBannerContainer}>
            running unsecure instance — authentication is disabled
        </div>
    )
}
