import cls from "@/pages/workbench/components/WorkbenchDangerZone/components/DangerZoneText.module.css"

export default function DangerZoneText() {
    return (
        <div className={cls.DangerZoneTextContainer}>
            <div className={cls.DangerZoneTitle}>Delete this Workbench</div>
            <div className={cls.DangerZoneSub}>
                Permanent. Destroys the container and its workspace — all files and state saved
                in it are lost. This is different from stopping it, which keeps everything intact.
            </div>
        </div>
    )
}
