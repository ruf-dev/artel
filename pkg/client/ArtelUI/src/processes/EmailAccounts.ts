import useUser from "@/hooks/user/User.ts"
import {EmailAccountInfo, EmailAccountsAPI, AddEmailAccountRequest} from "@/app/api/artel/email_accounts.pb.ts"

export interface IEmailAccountsService {
    list: () => Promise<EmailAccountInfo[]>
    add: (req: AddEmailAccountRequest) => Promise<void>
    remove: (id: string) => Promise<void>
}

export class EmailAccountsService implements IEmailAccountsService {
    async list(): Promise<EmailAccountInfo[]> {
        const res = await EmailAccountsAPI.ListEmailAccounts({}, useUser.getState().auth.getInitReq())
        return res.accounts ?? []
    }

    async add(req: AddEmailAccountRequest): Promise<void> {
        await EmailAccountsAPI.AddEmailAccount(req, useUser.getState().auth.getInitReq())
    }

    async remove(id: string): Promise<void> {
        await EmailAccountsAPI.DeleteEmailAccount({id}, useUser.getState().auth.getInitReq())
    }
}

export const emailAccountsService = new EmailAccountsService()
