import {CreateTriggerRequest, TractsAPI} from "@/app/api/artel/tracts.pb.ts"
import useUser from "@/hooks/user/User.ts"
import {
    CreatedTrigger,
    Tract,
    TractCondition,
    TractDefinition,
    TractRun,
    TractRunStep,
    TractTool,
    Trigger,
    TriggerSource,
    definitionToProto,
    toRun,
    toRunStep,
    toTool,
    toTract,
    toTrigger,
    toTriggerSource,
} from "@/processes/tractsMappers.ts"

export * from "@/processes/tractsMappers.ts"

export interface ITractsService {
    listTracts: () => Promise<Tract[]>
    getTract: (uuid: string) => Promise<Tract>
    createTract: (name: string, description: string, definition: TractDefinition) => Promise<{ tract: Tract; warnings: string[] }>
    updateTract: (uuid: string, name: string, description: string, definition: TractDefinition) => Promise<{ tract: Tract; warnings: string[] }>
    deleteTract: (uuid: string) => Promise<void>
    setTractEnabled: (uuid: string, enabled: boolean) => Promise<void>
    runTract: (tractUuid: string, params: unknown) => Promise<void>

    listRuns: (tractUuid: string, limit: number) => Promise<TractRun[]>
    getRun: (runUuid: string) => Promise<{ run: TractRun; steps: TractRunStep[] }>

    listTractTools: () => Promise<TractTool[]>
    listTriggerSources: () => Promise<TriggerSource[]>

    createTrigger: (name: string, kind: string, source: string, config: unknown, payloadSchema: unknown) => Promise<CreatedTrigger>
    listTriggers: () => Promise<Trigger[]>
    deleteTrigger: (uuid: string) => Promise<void>
    setTriggerEnabled: (uuid: string, enabled: boolean) => Promise<void>
    rotateTriggerToken: (uuid: string) => Promise<CreatedTrigger>
    linkTrigger: (triggerUuid: string, tractUuid: string, filters: TractCondition[]) => Promise<void>
    unlinkTrigger: (triggerUuid: string, tractUuid: string) => Promise<void>
}

export class TractsService implements ITractsService {
    async listTracts(): Promise<Tract[]> {
        const res = await TractsAPI.ListTracts({}, useUser.getState().auth.getInitReq())
        return (res.tracts ?? []).map(toTract)
    }

    async getTract(uuid: string): Promise<Tract> {
        const res = await TractsAPI.GetTract({uuid}, useUser.getState().auth.getInitReq())
        return toTract(res.tract!)
    }

    async createTract(name: string, description: string, definition: TractDefinition): Promise<{ tract: Tract; warnings: string[] }> {
        const res = await TractsAPI.CreateTract(
            {name, description, definition: definitionToProto(definition)},
            useUser.getState().auth.getInitReq(),
        )
        return {tract: toTract(res.tract!), warnings: res.warnings ?? []}
    }

    async updateTract(uuid: string, name: string, description: string, definition: TractDefinition): Promise<{ tract: Tract; warnings: string[] }> {
        const res = await TractsAPI.UpdateTract(
            {uuid, name, description, definition: definitionToProto(definition)},
            useUser.getState().auth.getInitReq(),
        )
        return {tract: toTract(res.tract!), warnings: res.warnings ?? []}
    }

    async deleteTract(uuid: string): Promise<void> {
        await TractsAPI.DeleteTract({uuid}, useUser.getState().auth.getInitReq())
    }

    async setTractEnabled(uuid: string, enabled: boolean): Promise<void> {
        await TractsAPI.SetTractEnabled({uuid, enabled}, useUser.getState().auth.getInitReq())
    }

    async runTract(tractUuid: string, params: unknown): Promise<void> {
        await TractsAPI.RunTract(
            {tractUuid, params: JSON.stringify(params ?? {})},
            useUser.getState().auth.getInitReq(),
        )
    }

    async listRuns(tractUuid: string, limit: number): Promise<TractRun[]> {
        const res = await TractsAPI.ListRuns({tractUuid, limit}, useUser.getState().auth.getInitReq())
        return (res.runs ?? []).map(toRun)
    }

    async getRun(runUuid: string): Promise<{ run: TractRun; steps: TractRunStep[] }> {
        const res = await TractsAPI.GetRun({runUuid}, useUser.getState().auth.getInitReq())
        return {run: toRun(res.run!), steps: (res.steps ?? []).map(toRunStep)}
    }

    async listTractTools(): Promise<TractTool[]> {
        const res = await TractsAPI.ListTractTools({}, useUser.getState().auth.getInitReq())
        return (res.tools ?? []).map(toTool)
    }

    async listTriggerSources(): Promise<TriggerSource[]> {
        const res = await TractsAPI.ListTriggerSources({}, useUser.getState().auth.getInitReq())
        return (res.sources ?? []).map(toTriggerSource)
    }

    async createTrigger(name: string, kind: string, source: string, config: unknown, payloadSchema: unknown): Promise<CreatedTrigger> {
        const req: CreateTriggerRequest = {
            name,
            kind,
            source,
            config: JSON.stringify(config ?? {}),
            payloadSchema: JSON.stringify(payloadSchema ?? {}),
        }
        const res = await TractsAPI.CreateTrigger(req, useUser.getState().auth.getInitReq())
        return {trigger: toTrigger(res.trigger!), webhookUrl: res.webhookUrl ?? "", webhookToken: res.webhookToken ?? ""}
    }

    async listTriggers(): Promise<Trigger[]> {
        const res = await TractsAPI.ListTriggers({}, useUser.getState().auth.getInitReq())
        return (res.triggers ?? []).map(toTrigger)
    }

    async deleteTrigger(uuid: string): Promise<void> {
        await TractsAPI.DeleteTrigger({uuid}, useUser.getState().auth.getInitReq())
    }

    async setTriggerEnabled(uuid: string, enabled: boolean): Promise<void> {
        await TractsAPI.SetTriggerEnabled({uuid, enabled}, useUser.getState().auth.getInitReq())
    }

    async rotateTriggerToken(uuid: string): Promise<CreatedTrigger> {
        const res = await TractsAPI.RotateTriggerToken({uuid}, useUser.getState().auth.getInitReq())
        return {trigger: toTrigger(res.trigger!), webhookUrl: res.webhookUrl ?? "", webhookToken: res.webhookToken ?? ""}
    }

    async linkTrigger(triggerUuid: string, tractUuid: string, filters: TractCondition[]): Promise<void> {
        await TractsAPI.LinkTrigger(
            {triggerUuid, tractUuid, filters: JSON.stringify(filters ?? [])},
            useUser.getState().auth.getInitReq(),
        )
    }

    async unlinkTrigger(triggerUuid: string, tractUuid: string): Promise<void> {
        await TractsAPI.UnlinkTrigger({triggerUuid, tractUuid}, useUser.getState().auth.getInitReq())
    }
}

export const tractsService = new TractsService()
