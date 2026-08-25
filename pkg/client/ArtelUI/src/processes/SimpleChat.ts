import useUser from "@/hooks/user/User.ts"
import {
    SimpleChat as SimpleChatProto,
    SimpleChatAPI,
    SimpleChatMessage as SimpleChatMessageProto,
} from "@/app/api/artel/simple_chat.pb.ts"

export interface SimpleChat {
    id: string
    vaultId: string
    title: string
    model: string
    vaultAccess: boolean
    createdAt: string
    updatedAt: string
    lastActivityAt: string
}

export interface SimpleChatMessage {
    id: string
    role: string
    content: string
    toolName?: string
    isError?: boolean
    model?: string
    seq: string
    createdAt: string
}

function toSimpleChat(c: SimpleChatProto): SimpleChat {
    return {
        id: c.id ?? "",
        vaultId: c.vaultId ?? "",
        title: c.title ?? "",
        model: c.model ?? "",
        vaultAccess: c.vaultAccess ?? false,
        createdAt: c.createdAt ?? "",
        updatedAt: c.updatedAt ?? "",
        lastActivityAt: c.lastActivityAt ?? "",
    }
}

function toSimpleChatMessage(m: SimpleChatMessageProto): SimpleChatMessage {
    return {
        id: m.id ?? "",
        role: m.role ?? "",
        content: m.content ?? "",
        toolName: m.toolName,
        isError: m.isError,
        model: m.model,
        seq: m.seq ?? "",
        createdAt: m.createdAt ?? "",
    }
}

export interface ISimpleChatService {
    createChat: (vaultId: string, model: string, vaultAccess: boolean) => Promise<SimpleChat>
    listChats: (vaultId: string) => Promise<SimpleChat[]>
    getChat: (chatId: string) => Promise<{ chat: SimpleChat; messages: SimpleChatMessage[] }>
    deleteChat: (chatId: string) => Promise<void>
}

// Thin wrapper over the generated simple_chat.pb.ts client, reshaping proto
// responses into plain TS objects — mirrors processes/Workbench.ts's shape.
export class SimpleChatService implements ISimpleChatService {
    async createChat(vaultId: string, model: string, vaultAccess: boolean): Promise<SimpleChat> {
        // `model` is passed through even though CreateSimpleChatRequest may not
        // declare it yet at the time this file was written — the backend track is
        // adding it in parallel (see processes/SimpleChat.ts's colocated note in the
        // task brief). Re-run `bun gen` before typechecking to pick up the field.
        const res = await SimpleChatAPI.CreateSimpleChat(
            {vaultId, model, vaultAccess},
            useUser.getState().auth.getInitReq(),
        )
        return toSimpleChat(res.chat!)
    }

    async listChats(vaultId: string): Promise<SimpleChat[]> {
        const res = await SimpleChatAPI.ListSimpleChats({vaultId}, useUser.getState().auth.getInitReq())
        return (res.chats ?? []).map(toSimpleChat)
    }

    async getChat(chatId: string): Promise<{ chat: SimpleChat; messages: SimpleChatMessage[] }> {
        const res = await SimpleChatAPI.GetSimpleChat({chatId}, useUser.getState().auth.getInitReq())
        return {
            chat: toSimpleChat(res.chat!),
            messages: (res.messages ?? []).map(toSimpleChatMessage),
        }
    }

    async deleteChat(chatId: string): Promise<void> {
        await SimpleChatAPI.DeleteSimpleChat({chatId}, useUser.getState().auth.getInitReq())
    }
}

export const simpleChatService = new SimpleChatService()
