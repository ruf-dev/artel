import {useCallback, useEffect, useRef, useState} from "react"
import type {KeyboardEvent} from "react"
import {useNavigate} from "react-router-dom"

import cls from "@/pages/home/HomePage.module.css"

import {AuthMiddleware} from "@/processes/AuthMiddleware.ts"
import {VaultService} from "@/processes/Vaults.ts"
import {Path} from "@/app/routing/Router.tsx"
import {VaultsAPI} from "@/app/api/artel/vaults.pb.ts"
import {VaultItem} from "@/app/api/artel/vaults.pb.ts"
import VaultCard from "@/pages/home/VaultCard.tsx"
import Topbar from "@/components/Topbar/Topbar.tsx"

interface Props {
    auth: AuthMiddleware
    onLogout: () => void
}

const vaultService = new VaultService()

export default function HomePage({auth, onLogout}: Props) {
    const navigate = useNavigate()
    const [dialogOpen, setDialogOpen] = useState(false)
    const [vaultName, setVaultName] = useState("")
    const [creating, setCreating] = useState(false)
    const [vaults, setVaults] = useState<VaultItem[]>([])
    const [loading, setLoading] = useState(true)
    const inputRef = useRef<HTMLInputElement>(null)

    useEffect(() => {
        if (!auth.isAuthenticated()) {
            navigate(Path.InitPage)
        }
    }, [auth, navigate])

    useEffect(() => {
        if (dialogOpen) {
            inputRef.current?.focus()
        }
    }, [dialogOpen])

    const fetchVaults = useCallback(async () => {
        setLoading(true)
        try {
            const list = await vaultService.ListVaults()
            setVaults(list)
        } finally {
            setLoading(false)
        }
    }, [])

    useEffect(() => {
        if (auth.isAuthenticated()) {
            void fetchVaults()
        }
    }, [auth, fetchVaults])

    function handleLogout() {
        onLogout()
        navigate(Path.InitPage)
    }

    function openDialog() {
        setVaultName("")
        setDialogOpen(true)
    }

    function closeDialog() {
        setDialogOpen(false)
    }

    async function handleCreate() {
        const name = vaultName.trim()
        if (!name) return
        setCreating(true)
        try {
            await VaultsAPI.CreateVault({name}, auth.getInitReq())
            closeDialog()
            void fetchVaults()
        } finally {
            setCreating(false)
        }
    }

    function handleKeyDown(e: KeyboardEvent<HTMLInputElement>) {
        if (e.key === "Enter") void handleCreate()
        if (e.key === "Escape") closeDialog()
    }

    return (
        <div className={cls.Root}>
            <Topbar onLogout={handleLogout}/>

            <div className={cls.Hero}>
                <div className={cls.HeroTitles}>
                    <div className={cls.Eyebrow}>Workspace</div>
                    <h1 className={cls.HeroTitle}>Your vaults</h1>
                    <p className={cls.HeroSub}>
                        <b>{loading ? "…" : `${vaults.length} ${vaults.length === 1 ? "vault" : "vaults"}`}</b>
                        {" · "}<span>all systems operational</span>
                    </p>
                </div>
                <button className={cls.NewVaultBtn} onClick={openDialog}>
                    <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                        <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
                    </svg>
                    New vault
                </button>
            </div>

            <div className={cls.Content}>
                {loading ? (
                    <p className={cls.Empty}>Loading…</p>
                ) : vaults.length === 0 ? (
                    <p className={cls.Empty}>No vaults yet. Create one with +</p>
                ) : (
                    <div className={cls.Grid}>
                        {vaults.map(v => (
                            <VaultCard key={v.id} vault={v}/>
                        ))}
                    </div>
                )}
            </div>

            {dialogOpen && (
                <div className={cls.Overlay} onClick={closeDialog}>
                    <div className={cls.Dialog} onClick={e => e.stopPropagation()}>
                        <h2 className={cls.DialogTitle}>New vault</h2>
                        <input
                            ref={inputRef}
                            className={cls.DialogInput}
                            placeholder="Vault name"
                            value={vaultName}
                            onChange={e => setVaultName(e.target.value)}
                            onKeyDown={handleKeyDown}
                            disabled={creating}
                        />
                        <div className={cls.DialogActions}>
                            <button className={cls.CancelBtn} onClick={closeDialog} disabled={creating}>Cancel</button>
                            <button className={cls.CreateBtn} onClick={handleCreate} disabled={creating || !vaultName.trim()}>
                                {creating ? "Creating…" : "Create"}
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    )
}
