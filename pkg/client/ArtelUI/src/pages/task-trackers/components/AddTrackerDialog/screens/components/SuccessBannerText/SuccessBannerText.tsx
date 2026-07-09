export default function SuccessBannerText({boardCount}: {boardCount: number}) {
    return (
        <span>
            Successfully connected.{" "}
            <b>{boardCount} {boardCount === 1 ? "board" : "boards"}</b> available.
        </span>
    )
}
