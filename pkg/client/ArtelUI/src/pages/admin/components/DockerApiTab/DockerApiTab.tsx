import {useState, useEffect, useCallback, useRef} from "react"

import cls from "@/pages/admin/components/DockerApiTab/DockerApiTab.module.css"
import {DockerHostsAPI, GetDockerHostResponse} from "@/app/api/artel/docker_hosts.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import useUser from "@/hooks/user/User.ts"
import DockerHostsActionBar
    from "@/pages/admin/components/DockerApiTab/components/DockerHostsActionBar/DockerHostsActionBar.tsx"
import DockerHostList from "@/pages/admin/components/DockerApiTab/components/DockerHostList/DockerHostList.tsx"
import DockerHostFormDialog, {
    DockerHostFormDialogSaveData,
} from "@/pages/admin/components/DockerApiTab/components/DockerHostFormDialog/DockerHostFormDialog.tsx"

// TLS fields are write-only and optional on the wire (proto3 `optional`): an
// empty textarea must be omitted from the request, not sent as "", so the
// backend keeps the previously stored value on edit instead of clearing it.
function tlsPatch(data: DockerHostFormDialogSaveData) {
    return {
        caCert: data.caCert || undefined,
        clientCert: data.clientCert || undefined,
        clientKey: data.clientKey || undefined,
    }
}

interface Props {
    autoOpenAddDialog?: boolean
}

export default function DockerApiTab({autoOpenAddDialog}: Props) {
    const {auth} = useUser()
    const {OpenDialog} = useDialog()
    const [hosts, setHosts] = useState<GetDockerHostResponse[]>([])
    const [loading, setLoading] = useState(true)
    const autoOpenHandled = useRef(false)

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
                onSave={async (data) => {
                    await DockerHostsAPI.RegisterDockerHost(
                        {url: data.url, ...tlsPatch(data)},
                        auth.getInitReq()
                    )
                    await loadHosts()
                }}
            />
        )
    }

    function openEditDialog(host: GetDockerHostResponse) {
        OpenDialog(
            <DockerHostFormDialog
                initial={host}
                onSave={async (data) => {
                    await DockerHostsAPI.UpdateDockerHost(
                        {id: host.id, url: data.url, ...tlsPatch(data)},
                        auth.getInitReq()
                    )
                    await loadHosts()
                }}
            />
        )
    }

    useEffect(() => {
        if (autoOpenAddDialog && !autoOpenHandled.current) {
            autoOpenHandled.current = true
            openAddDialog()
        }
    }, [autoOpenAddDialog])

    return (
        <div className={cls.DockerApiTabContainer}>
            <DockerHostsActionBar count={hosts.length} onAddClick={openAddDialog} />
            <DockerHostList hosts={hosts} loading={loading} onEdit={openEditDialog} onDelete={handleDelete} />
        </div>
    )
}
