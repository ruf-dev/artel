import {create} from 'zustand'

import {EmailAccountInfo, AddEmailAccountRequest} from "@/app/api/artel/email_accounts.pb.ts"
import {emailAccountsService} from "@/processes/EmailAccounts.ts"

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
            const accounts = await emailAccountsService.list()
            set({accounts})
        } finally {
            set({loading: false})
        }
    },

    add: async (req: AddEmailAccountRequest) => {
        await emailAccountsService.add(req)
        await get().fetch()
    },

    remove: async (id: string) => {
        await emailAccountsService.remove(id)
        await get().fetch()
    },
}))
