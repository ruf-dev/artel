import cls from "@/pages/admin/components/DockerApiTab/DockerApiTab.module.css"

export default function DockerApiTab() {
    return (
        <div className={cls.DockerApiTabContainer}>
            <p className={cls.EmptyState}>
                Docker API connection management is coming soon. This will let you register
                the Docker host(s) used to run per-vault Workbench containers.
            </p>
        </div>
    )
}
