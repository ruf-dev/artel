import cls from "@/pages/error/ErrorPage.module.css"

export default function ErrorPage() {
    return (
        <div className={cls.ErrorContainer}>
            <h1>Something went wrong</h1>
            <a href="/" className={cls.HomeLink}>Go home</a>
        </div>
    )
}
