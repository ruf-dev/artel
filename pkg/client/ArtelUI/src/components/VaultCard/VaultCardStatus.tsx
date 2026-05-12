import cls from "@/components/VaultCard/VaultCardStatus.module.css"

export default function VaultCardStatus() {
    return (
        <div className={cls.VaultCardStatusContainer}>
            <span className={cls.Chip}>
                <span className={cls.StatusDot} />
                <span>Online</span>
            </span>
        </div>
    )
}
