import {GetDockerHostResponse} from "@/app/api/artel/docker_hosts.pb.ts"
import DockerHostRow from "@/pages/admin/components/DockerApiTab/components/DockerHostRow/DockerHostRow.tsx"
import cls from "@/pages/admin/components/DockerApiTab/components/DockerHostList/DockerHostList.module.css"

interface DockerHostListProps {
    hosts: GetDockerHostResponse[]
    loading: boolean
    onEdit: (host: GetDockerHostResponse) => void
    onDelete: (id: string) => Promise<void>
}

export default function DockerHostList({hosts, loading, onEdit, onDelete}: DockerHostListProps) {
    if (loading) {
        return <p className={cls.Empty}>Loading…</p>
    }
    if (hosts.length === 0) {
        return <p className={cls.Empty}>No Docker hosts yet. Add one to get started.</p>
    }
    return (
        <>
            {hosts.map(host => (
                <DockerHostRow key={host.id} host={host} onEdit={onEdit} onDelete={onDelete} />
            ))}
        </>
    )
}
