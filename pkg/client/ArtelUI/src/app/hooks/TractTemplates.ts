import {create} from 'zustand'

import {TractTemplate, TractTemplateSummary, tractsService} from "@/processes/Tracts.ts"

interface TractTemplatesState {
    templates: TractTemplateSummary[]
    loading: boolean
    mineOnly: boolean
    category: string

    fetch: () => Promise<void>
    setMineOnly: (mineOnly: boolean) => void
    setCategory: (category: string) => void
    publish: (tractUuid: string, category: string) => Promise<TractTemplate>
    unpublish: (uuid: string) => Promise<void>
}

export const useTractTemplates = create<TractTemplatesState>((set, get) => ({
    templates: [],
    loading: false,
    mineOnly: false,
    category: "",

    fetch: async () => {
        set({loading: true})
        try {
            const templates = await tractsService.listTemplates(get().category, get().mineOnly)
            set({templates})
        } finally {
            set({loading: false})
        }
    },

    setMineOnly: (mineOnly: boolean) => {
        set({mineOnly})
        void get().fetch()
    },

    setCategory: (category: string) => {
        set({category})
        void get().fetch()
    },

    publish: async (tractUuid: string, category: string) => {
        const template = await tractsService.publishTemplate(tractUuid, category)
        await get().fetch()
        return template
    },

    unpublish: async (uuid: string) => {
        await tractsService.unpublishTemplate(uuid)
        await get().fetch()
    },
}))
