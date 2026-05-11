export default function ErrorPage() {
    return (
        <div style={{
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            height: "100%",
            flexDirection: "column",
            gap: "1rem",
        }}>
            <h1>Something went wrong</h1>
            <a href="/" style={{color: "var(--accent-fg-color)"}}>Go home</a>
        </div>
    )
}
