import {Button} from "@vervstack/chures"

import cls from "@/components/EmptyState/EmptyState.module.css"

interface Props {
    onCreateClick: () => void
}

export default function EmptyState({onCreateClick}: Props) {
    return (
        <div className={cls.Root}>
            <div className={cls.Art} aria-hidden="true">
                <svg
                    viewBox="0 0 100 100" fill="none" stroke="var(--color-empty-icon-stroke)"
                    strokeWidth="1.4" strokeLinecap="round" strokeLinejoin="round"
                >
                    <rect x="18" y="28" width="64" height="52" rx="8"/>
                    <path d="M30 28v-6a8 8 0 0 1 8-8h24a8 8 0 0 1 8 8v6"/>
                    <circle cx="50" cy="54" r="6" fill="rgba(255,75,62,0.18)" stroke="#FF4B3E"/>
                    <line x1="50" y1="60" x2="50" y2="68"/>
                </svg>
            </div>
            <h2 className={cls.Title}>No vaults yet</h2>
            <p className={cls.Sub}>
                Vaults are encrypted CouchDB instances for your data. Create your first one
                and you'll get a connection string in seconds.
            </p>
            <Button variant="primary" onClick={onCreateClick}>
                <svg
                    viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor"
                    strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"
                >
                    <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                Create your first vault
            </Button>
        </div>
    )
}
