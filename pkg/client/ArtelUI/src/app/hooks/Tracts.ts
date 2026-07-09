import {create} from 'zustand'
import {useQuery} from '@tanstack/react-query'

import {CreatedTrigger, Tract, TractDefinition, TractRun, TractRunStep, TractTool, tractsService, Trigger} from "@/processes/Tracts.ts"
import {retryOnStatus} from "@/processes/grpcErrors.ts"
import useUser from "@/hooks/user/User.ts"

export const triggerSourcesQueryKey = ['triggerSources'] as const

export function useTriggerSources() {
    const {auth} = useUser()

    const q = useQuery({
        queryKey: triggerSourcesQueryKey,
        queryFn: () => tractsService.listTriggerSources(),
        enabled: auth.isAuthenticated(),
        retry: retryOnStatus(),
    })

    return {
        triggerSources: q.data ?? [],
        isLoading: q.isLoading,
        error: q.error,
    }
}

export const triggersQueryKey = ['triggers'] as const

export function useTrigger(uuid: string) {
    const {auth} = useUser()

    const q = useQuery({
        queryKey: triggersQueryKey,
        queryFn: () => tractsService.listTriggers(),
        enabled: auth.isAuthenticated() && !!uuid,
        select: (triggers) => triggers.find(t => t.uuid === uuid),
        retry: retryOnStatus(),
    })

    return {
        trigger: q.data,
        isLoading: q.isLoading,
        error: q.error,
    }
}

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

    fetchTriggers: () => Promise<void>
    createAndLinkTrigger: (name: string, kind: string, source: string, payloadSchema: unknown, tractUuid: string) => Promise<CreatedTrigger>
    deleteTrigger: (uuid: string) => Promise<void>
    setTriggerEnabled: (uuid: string, enabled: boolean) => Promise<void>
    rotateTriggerToken: (uuid: string) => Promise<{ trigger: Trigger; webhookUrl: string; webhookToken: string }>
    unlinkTrigger: (triggerUuid: string, tractUuid: string) => Promise<void>
}

async function handleSetTriggerEnabled(uuid: string, enabled: boolean, get: () => TractsState) {
    await tractsService.setTriggerEnabled(uuid, enabled)
    await get().fetchTriggers()
}

async function handleRotateTriggerToken(uuid: string, get: () => TractsState) {
    const result = await tractsService.rotateTriggerToken(uuid)
    await get().fetchTriggers()
    return result
}

async function handleCreateAndLinkTrigger(name: string, kind: string, source: string, payloadSchema: unknown, tractUuid: string, get: () => TractsState) {
    const result = await tractsService.createTrigger(name, kind, source, {}, payloadSchema)
    await tractsService.linkTrigger(result.trigger.uuid, tractUuid, [])
    await Promise.all([get().fetchTriggers(), get().fetch()])
    return result
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

    fetchTriggers: async () => {
        set({triggersLoading: true})
        try {
            const triggers = await tractsService.listTriggers()
            set({triggers})
        } finally {
            set({triggersLoading: false})
        }
    },

    createAndLinkTrigger: async (name: string, kind: string, source: string, payloadSchema: unknown, tractUuid: string) => handleCreateAndLinkTrigger(name, kind, source, payloadSchema, tractUuid, get),

    deleteTrigger: async (uuid: string) => {
        await tractsService.deleteTrigger(uuid)
        await get().fetchTriggers()
        await get().fetch()
    },

    setTriggerEnabled: async (uuid: string, enabled: boolean) => handleSetTriggerEnabled(uuid, enabled, get),

    rotateTriggerToken: async (uuid: string) => handleRotateTriggerToken(uuid, get),

    unlinkTrigger: async (triggerUuid: string, tractUuid: string) => {
        await tractsService.unlinkTrigger(triggerUuid, tractUuid)
        await get().fetch()
    },
}))
