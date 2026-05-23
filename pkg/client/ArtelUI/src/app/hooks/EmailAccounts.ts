import {create} from 'zustand'

import {EmailAccountInfo, EmailAccountsAPI, AddEmailAccountRequest} from "@/app/api/artel/email_accounts.pb.ts"
import useUser from "@/hooks/user/User.ts"

interface EmailAccountsState {
    accounts: EmailAccountInfo[]
    loading: boolean
    fetch: () => Promise<void>
    add: (req: AddEmailAccountRequest) => Promise<void>
    remove: (id: string) => Promise<void>
}

export const useEmailAccounts = create<EmailAccountsState>((set, get) => ({
    accounts: [],
    loading: false,

    fetch: async () => {
        set({loading: true})
        try {
            const {auth} = useUser.getState()
            const resp = await EmailAccountsAPI.ListEmailAccounts({}, auth.getInitReq())
            set({accounts: resp.accounts ?? []})
        } finally {
            set({loading: false})
        }
    },

    add: async (req: AddEmailAccountRequest) => {
        const {auth} = useUser.getState()
        await EmailAccountsAPI.AddEmailAccount(req, auth.getInitReq())
        await get().fetch()
    },

    remove: async (id: string) => {
        const {auth} = useUser.getState()
        await EmailAccountsAPI.DeleteEmailAccount({id}, auth.getInitReq())
        await get().fetch()
    },
}))
