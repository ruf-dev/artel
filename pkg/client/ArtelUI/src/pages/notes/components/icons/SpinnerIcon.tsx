import cls from "@/pages/notes/components/SaveStatusIndicator/SaveStatusIndicator.module.css"

export default function SpinnerIcon() {
    return (
        <svg viewBox="0 0 12 12" width="10" height="10" className={cls.Spinner}>
            <circle cx="6" cy="6" r="4.5" fill="none" stroke="rgba(255,255,255,0.18)" strokeWidth="1.5" />
            <path d="M6 1.5 A4.5 4.5 0 0 1 10.5 6" fill="none" stroke="rgba(255,255,255,0.55)" strokeWidth="1.5" strokeLinecap="round" />
        </svg>
    )
}
