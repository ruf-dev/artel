import cls from "@/components/VaultCard/VaultCardQR.module.css"

export default function VaultCardQR() {
    return (
        <div className={cls.VaultCardQrContainer}>
            <div className={cls.QRPlaceholder} />
        </div>
    )
}
