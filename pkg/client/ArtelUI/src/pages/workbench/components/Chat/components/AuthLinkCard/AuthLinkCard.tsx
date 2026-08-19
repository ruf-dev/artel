import cls from "@/pages/workbench/components/Chat/components/AuthLinkCard/AuthLinkCard.module.css"

interface Props {
    url: string
}

export default function AuthLinkCard({url}: Props) {
    return (
        <div className={cls.AuthLinkCardContainer}>
            <p className={cls.Text}>Claude needs you to authorize access to continue.</p>
            <a className={cls.LinkButton} href={url} target="_blank" rel="noreferrer">Authorize Claude</a>
        </div>
    )
}
