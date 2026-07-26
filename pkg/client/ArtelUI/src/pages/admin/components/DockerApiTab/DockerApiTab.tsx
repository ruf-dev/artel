import {useState, useEffect, useCallback} from "react"

import cls from "@/pages/admin/components/DockerApiTab/DockerApiTab.module.css"
import {DockerHostsAPI, GetDockerHostResponse} from "@/app/api/artel/docker_hosts.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import useUser from "@/hooks/user/User.ts"
import DockerHostsActionBar
    from "@/pages/admin/components/DockerApiTab/components/DockerHostsActionBar/DockerHostsActionBar.tsx"
import DockerHostList from "@/pages/admin/components/DockerApiTab/components/DockerHostList/DockerHostList.tsx"
import DockerHostFormDialog
    from "@/pages/admin/components/DockerApiTab/components/DockerHostFormDialog/DockerHostFormDialog.tsx"

export default function DockerApiTab() {
    const {auth} = useUser()
    const {OpenDialog} = useDialog()
    const [hosts, setHosts] = useState<GetDockerHostResponse[]>([])
    const [loading, setLoading] = useState(true)

    const loadHosts = useCallback(async () => {
        setLoading(true)
        try {
            const res = await DockerHostsAPI.ListDockerHosts({}, auth.getInitReq())
            setHosts(res.hosts ?? [])
        } finally {
            setLoading(false)
        }
    }, [auth])

    useEffect(() => { void loadHosts() }, [loadHosts])

    async function handleDelete(id: string) {
        await DockerHostsAPI.DeleteDockerHost({id}, auth.getInitReq())
        await loadHosts()
    }

    function openAddDialog() {
        OpenDialog(
            <DockerHostFormDialog
                onSave={async (url) => {
                    await DockerHostsAPI.RegisterDockerHost({url}, auth.getInitReq())
                    await loadHosts()
                }}
            />
        )
    }

    function openEditDialog(host: GetDockerHostResponse) {
        OpenDialog(
            <DockerHostFormDialog
                initial={host}
                onSave={async (url) => {
                    await DockerHostsAPI.UpdateDockerHost(
                        {id: host.id, url},
                        auth.getInitReq()
                    )
                    await loadHosts()
                }}
            />
        )
    }

    return (
        <div className={cls.DockerApiTabContainer}>
            <DockerHostsActionBar count={hosts.length} onAddClick={openAddDialog} />
            <DockerHostList hosts={hosts} loading={loading} onEdit={openEditDialog} onDelete={handleDelete} />
        </div>
    )
}
