import {create} from 'zustand'

import {SkillInfo} from "@/app/api/artel/skills.pb.ts"
import {skillsService} from "@/processes/Skills.ts"

interface SkillsState {
    skills: SkillInfo[]
    loading: boolean
    fetch: (vaultId: string) => Promise<void>
    create: (
        vaultId: string, name: string, description: string, storageMode: string, body: string, hotPlug: boolean,
    ) => Promise<SkillInfo>
    update: (
        vaultId: string, slug: string, name: string, description: string, storageMode: string, body: string,
    ) => Promise<SkillInfo>
    remove: (vaultId: string, slug: string) => Promise<void>
    setHotPlug: (vaultId: string, slug: string, hotPlug: boolean) => Promise<void>
}

export const useSkills = create<SkillsState>((set, get) => ({
    skills: [],
    loading: false,

    fetch: async (vaultId: string) => {
        set({loading: true})
        try {
            const skills = await skillsService.list(vaultId)
            set({skills})
        } finally {
            set({loading: false})
        }
    },

    create: async (
        vaultId: string, name: string, description: string, storageMode: string, body: string, hotPlug: boolean,
    ) => {
        const skill = await skillsService.create(vaultId, name, description, storageMode, body, hotPlug)
        await get().fetch(vaultId)
        return skill
    },

    update: async (
        vaultId: string, slug: string, name: string, description: string, storageMode: string, body: string,
    ) => {
        const skill = await skillsService.update(vaultId, slug, name, description, storageMode, body)
        await get().fetch(vaultId)
        return skill
    },

    remove: async (vaultId: string, slug: string) => {
        await skillsService.remove(vaultId, slug)
        await get().fetch(vaultId)
    },

    setHotPlug: async (vaultId: string, slug: string, hotPlug: boolean) => {
        await skillsService.setHotPlug(vaultId, slug, hotPlug)
        await get().fetch(vaultId)
    },
}))
