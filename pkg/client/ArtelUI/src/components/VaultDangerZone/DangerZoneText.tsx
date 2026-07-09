import cls from "@/components/VaultDangerZone/DangerZoneText.module.css"

export default function DangerZoneText() {
    return (
        <div className={cls.DangerZoneTextContainer}>
            <div className={cls.DangerZoneTitle}>Delete this vault</div>
            <div className={cls.DangerZoneSub}>Permanent. Connection string stops working immediately.</div>
        </div>
    )
}
