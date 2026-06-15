import cls from "@/pages/service-unavailable/ServiceUnavailablePage.module.css"

export default function ServiceUnavailablePage() {
    return (
        <div className={cls.ServiceUnavailableContainer}>
            <div className={cls.Card}>
                <div className={cls.Logo}>artel</div>
                <div className={cls.Body}>
                    <div className={cls.PulsingDot}/>
                    <h2 className={cls.Title}>Service not responding</h2>
                    <p className={cls.Sub}>Cannot reach the server. Retrying automatically...</p>
                </div>
            </div>
        </div>
    )
}
