// S3/MinIO has no bundled brand SVG yet — mirrors the emoji used by the BYOK
// "Infrastructure" card before it became a real connection (BYOKSection.tsx).
export default function S3Icon() {
    return <span aria-hidden="true" style={{fontSize: "1.25rem", lineHeight: 1}}>🪣</span>
}
