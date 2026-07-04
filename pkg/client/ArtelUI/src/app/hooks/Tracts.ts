import {create} from 'zustand'

import {Tract, TractCondition, TractDefinition, TractRun, TractRunStep, TractTool, tractsService, Trigger, TriggerSource} from "@/processes/Tracts.ts"

interface TractsState {
    tracts: Tract[]
    loading: boolean

    runsByTract: Record<string, TractRun[]>
    runsLoading: boolean

    currentRun: TractRun | null
    currentRunSteps: TractRunStep[]
    currentRunLoading: boolean

    tools: TractTool[]
    toolsLoading: boolean
    triggerSources: TriggerSource[]
    triggerSourcesLoading: boolean

    triggers: Trigger[]
    triggersLoading: boolean

    fetch: () => Promise<void>
    createTract: (name: string, description: string, definition: TractDefinition) => Promise<{ tract: Tract; warnings: string[] }>
    updateTract: (uuid: string, name: string, description: string, definition: TractDefinition) => Promise<{ tract: Tract; warnings: string[] }>
    deleteTract: (uuid: string) => Promise<void>
    setEnabled: (uuid: string, enabled: boolean) => Promise<void>
    runTract: (tractUuid: string, params: unknown) => Promise<void>

    fetchRuns: (tractUuid: string) => Promise<void>
    fetchRun: (runUuid: string) => Promise<void>
    clearCurrentRun: () => void

    fetchTools: () => Promise<void>
    fetchTriggerSources: () => Promise<void>

    fetchTriggers: () => Promise<void>
    createTrigger: (name: string, kind: string, source: string, config: unknown, payloadSchema: unknown) => Promise<{ trigger: Trigger; webhookUrl: string; webhookToken: string }>
    deleteTrigger: (uuid: string) => Promise<void>
    setTriggerEnabled: (uuid: string, enabled: boolean) => Promise<void>
    rotateTriggerToken: (uuid: string) => Promise<{ trigger: Trigger; webhookUrl: string; webhookToken: string }>
    linkTrigger: (triggerUuid: string, tractUuid: string, filters: TractCondition[]) => Promise<void>
    unlinkTrigger: (triggerUuid: string, tractUuid: string) => Promise<void>
}

export const useTracts = create<TractsState>((set, get) => ({
    tracts: [],
    loading: false,

    runsByTract: {},
    runsLoading: false,

    currentRun: null,
    currentRunSteps: [],
    currentRunLoading: false,

    tools: [],
    toolsLoading: false,
    triggerSources: [],
    triggerSourcesLoading: false,

    triggers: [],
    triggersLoading: false,

    fetch: async () => {
        set({loading: true})
        try {
            const tracts = await tractsService.listTracts()
            set({tracts})
        } finally {
            set({loading: false})
        }
    },

    createTract: async (name: string, description: string, definition: TractDefinition) => {
        const {tract, warnings} = await tractsService.createTract(name, description, definition)
        await get().fetch()
        return {tract, warnings}
    },

    updateTract: async (uuid: string, name: string, description: string, definition: TractDefinition) => {
        const {tract, warnings} = await tractsService.updateTract(uuid, name, description, definition)
        await get().fetch()
        return {tract, warnings}
    },

    deleteTract: async (uuid: string) => {
        await tractsService.deleteTract(uuid)
        await get().fetch()
    },

    setEnabled: async (uuid: string, enabled: boolean) => {
        await tractsService.setTractEnabled(uuid, enabled)
        await get().fetch()
    },

    runTract: async (tractUuid: string, params: unknown) => {
        await tractsService.runTract(tractUuid, params)
        // RunTract's response is intentionally empty — refetch to see the new run.
        await get().fetchRuns(tractUuid)
    },

    fetchRuns: async (tractUuid: string) => {
        set({runsLoading: true})
        try {
            const runs = await tractsService.listRuns(tractUuid, 50)
            set({runsByTract: {...get().runsByTract, [tractUuid]: runs}})
        } finally {
            set({runsLoading: false})
        }
    },

    fetchRun: async (runUuid: string) => {
        set({currentRunLoading: true})
        try {
            const {run, steps} = await tractsService.getRun(runUuid)
            set({currentRun: run, currentRunSteps: steps})
        } finally {
            set({currentRunLoading: false})
        }
    },

    clearCurrentRun: () => set({currentRun: null, currentRunSteps: []}),

    fetchTools: async () => {
        set({toolsLoading: true})
        try {
            const tools = await tractsService.listTractTools()
            set({tools})
        } finally {
            set({toolsLoading: false})
        }
    },

    fetchTriggerSources: async () => {
        set({triggerSourcesLoading: true})
        try {
            const triggerSources = await tractsService.listTriggerSources()
            set({triggerSources})
        } finally {
            set({triggerSourcesLoading: false})
        }
    },

    fetchTriggers: async () => {
        set({triggersLoading: true})
        try {
            const triggers = await tractsService.listTriggers()
            set({triggers})
        } finally {
            set({triggersLoading: false})
        }
    },

    createTrigger: async (name: string, kind: string, source: string, config: unknown, payloadSchema: unknown) => {
        const result = await tractsService.createTrigger(name, kind, source, config, payloadSchema)
        await get().fetchTriggers()
        return result
    },

    deleteTrigger: async (uuid: string) => {
        await tractsService.deleteTrigger(uuid)
        await get().fetchTriggers()
        await get().fetch()
    },

    setTriggerEnabled: async (uuid: string, enabled: boolean) => {
        await tractsService.setTriggerEnabled(uuid, enabled)
        await get().fetchTriggers()
    },

    rotateTriggerToken: async (uuid: string) => {
        const result = await tractsService.rotateTriggerToken(uuid)
        await get().fetchTriggers()
        return result
    },

    linkTrigger: async (triggerUuid: string, tractUuid: string, filters: TractCondition[]) => {
        await tractsService.linkTrigger(triggerUuid, tractUuid, filters)
        await get().fetch()
    },

    unlinkTrigger: async (triggerUuid: string, tractUuid: string) => {
        await tractsService.unlinkTrigger(triggerUuid, tractUuid)
        await get().fetch()
    },
}))
