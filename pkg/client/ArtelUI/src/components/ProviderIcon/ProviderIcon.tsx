import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"

export default function ProviderIcon({provider}: {provider?: ExternalProvider}) {
    switch (provider) {
        case ExternalProvider.EXTERNAL_PROVIDER_GOOGLE_SHEETS:
            return <GoogleSheetsIcon/>
        case ExternalProvider.EXTERNAL_PROVIDER_TRELLO:
            return <TrelloIcon/>
        case ExternalProvider.EXTERNAL_PROVIDER_MIRO:
            return <MiroIcon/>
        case ExternalProvider.EXTERNAL_PROVIDER_EMAIL:
            return <EmailIcon/>
        case ExternalProvider.EXTERNAL_PROVIDER_GITLAB:
            return <GitlabIcon/>
        default:
            return <UnknownProviderIcon/>
    }
}

function GoogleSheetsIcon() {
    return (
        <svg width="20" height="20" viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
            <rect x="10" y="2" width="28" height="44" rx="2" fill="#23A566"/>
            <rect x="26" y="2" width="12" height="12" fill="#159C53"/>
            <path d="M26 2L38 14H26V2Z" fill="#82C9A7"/>
            <rect x="16" y="20" width="16" height="2" rx="1" fill="white" fillOpacity="0.9"/>
            <rect x="16" y="25" width="16" height="2" rx="1" fill="white" fillOpacity="0.9"/>
            <rect x="16" y="30" width="10" height="2" rx="1" fill="white" fillOpacity="0.9"/>
        </svg>
    )
}

function TrelloIcon() {
    return (
        <svg width="20" height="20" viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
            <rect width="48" height="48" rx="6" fill="#0052CC"/>
            <rect x="7" y="9" width="14" height="20" rx="2" fill="white"/>
            <rect x="27" y="9" width="14" height="13" rx="2" fill="white"/>
        </svg>
    )
}

function MiroIcon() {
    return (
        <svg width="20" height="20" viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
            <rect width="48" height="48" rx="6" fill="#FFD02F"/>
            <path d="M32 10H26L30 24L22 10H16L20 24L12 10H6L16 38H22L18 24L26 38H32L28 24L36 38H42L32 10Z" fill="#050038"/>
        </svg>
    )
}

function EmailIcon() {
    return (
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"
             stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
            <rect x="2" y="4" width="20" height="16" rx="2"/>
            <path d="M2 7l10 7 10-7"/>
        </svg>
    )
}

function GitlabIcon() {
    return (
        <svg aria-hidden="true" role="img" className="tanuki-logo" width="25" height="24" viewBox="0 0 25 24"
             fill="none" xmlns="http://www.w3.org/2000/svg">
            <path className="tanuki-shape tanuki"
                  d="m24.507 9.5-.034-.09L21.082.562a.896.896 0 0 0-1.694.091l-2.29 7.01H7.825L5.535.653a.898.898 0 0 0-1.694-.09L.451 9.411.416 9.5a6.297 6.297 0 0 0 2.09 7.278l.012.01.03.022 5.16 3.867 2.56 1.935 1.554 1.176a1.051 1.051 0 0 0 1.268 0l1.555-1.176 2.56-1.935 5.197-3.89.014-.01A6.297 6.297 0 0 0 24.507 9.5Z"
                  fill="#E24329"></path>
            <path className="tanuki-shape right-cheek"
                  d="m24.507 9.5-.034-.09a11.44 11.44 0 0 0-4.56 2.051l-7.447 5.632 4.742 3.584 5.197-3.89.014-.01A6.297 6.297 0 0 0 24.507 9.5Z"
                  fill="#FC6D26"></path>
            <path className="tanuki-shape chin"
                  d="m7.707 20.677 2.56 1.935 1.555 1.176a1.051 1.051 0 0 0 1.268 0l1.555-1.176 2.56-1.935-4.743-3.584-4.755 3.584Z"
                  fill="#FCA326"></path>
            <path className="tanuki-shape left-cheek"
                  d="M5.01 11.461a11.43 11.43 0 0 0-4.56-2.05L.416 9.5a6.297 6.297 0 0 0 2.09 7.278l.012.01.03.022 5.16 3.867 4.745-3.584-7.444-5.632Z"
                  fill="#FC6D26"></path>
        </svg>
    )
}

// TODO: placeholder glyph for providers without a dedicated brand icon yet - replace per-provider as icons are added.
function UnknownProviderIcon() {
    return (
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"
             stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
            <circle cx="12" cy="12" r="9"/>
            <path d="M9.5 9.5a2.5 2.5 0 1 1 3.5 2.3c-.7.35-1 .85-1 1.7v.3"/>
            <path d="M12 17h.01"/>
        </svg>
    )
}
