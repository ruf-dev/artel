import useUser from "@/hooks/user/User.ts"
import {
    SkillInfo,
    SkillsAPI,
} from "@/app/api/artel/skills.pb.ts"

export interface ISkillsService {
    list: (vaultId: string) => Promise<SkillInfo[]>
    get: (vaultId: string, slug: string) => Promise<{ skill: SkillInfo, body: string }>
    create: (
        vaultId: string, name: string, description: string, storageMode: string, body: string, hotPlug: boolean,
    ) => Promise<SkillInfo>
    update: (
        vaultId: string, slug: string, name: string, description: string, storageMode: string, body: string,
    ) => Promise<SkillInfo>
    remove: (vaultId: string, slug: string) => Promise<void>
    setHotPlug: (vaultId: string, slug: string, hotPlug: boolean) => Promise<SkillInfo>
}

export class SkillsService implements ISkillsService {
    async list(vaultId: string): Promise<SkillInfo[]> {
        const res = await SkillsAPI.ListSkills({vaultId}, useUser.getState().auth.getInitReq())
        return res.skills ?? []
    }

    async get(vaultId: string, slug: string): Promise<{ skill: SkillInfo, body: string }> {
        const res = await SkillsAPI.GetSkill({vaultId, slug}, useUser.getState().auth.getInitReq())
        return {skill: res.skill ?? {}, body: res.body ?? ""}
    }

    async create(
        vaultId: string, name: string, description: string, storageMode: string, body: string, hotPlug: boolean,
    ): Promise<SkillInfo> {
        const res = await SkillsAPI.CreateSkill(
            {vaultId, name, description, storageMode, body, hotPlug}, useUser.getState().auth.getInitReq(),
        )
        return res.skill ?? {}
    }

    async update(
        vaultId: string, slug: string, name: string, description: string, storageMode: string, body: string,
    ): Promise<SkillInfo> {
        const res = await SkillsAPI.UpdateSkill(
            {vaultId, slug, name, description, storageMode, body}, useUser.getState().auth.getInitReq(),
        )
        return res.skill ?? {}
    }

    async remove(vaultId: string, slug: string): Promise<void> {
        await SkillsAPI.DeleteSkill({vaultId, slug}, useUser.getState().auth.getInitReq())
    }

    async setHotPlug(vaultId: string, slug: string, hotPlug: boolean): Promise<SkillInfo> {
        const res = await SkillsAPI.SetSkillHotPlug(
            {vaultId, slug, hotPlug}, useUser.getState().auth.getInitReq(),
        )
        return res.skill ?? {}
    }
}

export const skillsService = new SkillsService()
