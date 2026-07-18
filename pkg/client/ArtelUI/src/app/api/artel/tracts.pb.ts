/* eslint-disable */
// @ts-nocheck

/**
 * This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
 */

import * as fm from "./fetch.pb";

type Absent<T, K extends keyof T> = { [k in Exclude<keyof T, K>]?: undefined };

type OneOf<T> =
  | { [k in keyof T]?: undefined }
  | (keyof T extends infer K
      ? K extends string & keyof T
        ? { [k in K]: T[K] } & Absent<T, K>
        : never
      : never);

export enum ScriptLanguage {
  SCRIPT_LANGUAGE_UNSPECIFIED = "SCRIPT_LANGUAGE_UNSPECIFIED",
  SCRIPT_LANGUAGE_JAVASCRIPT = "SCRIPT_LANGUAGE_JAVASCRIPT",
}

export type TractTriggerSummary = {
  uuid?: string;
  name?: string;
  kind?: string;
  source?: string;
};

export type TractLastRun = {
  status?: string;
  at?: string;
};

export type TractCondition = {
  left?: string;
  op?: string;
  right?: string;
};

export type ActionStep = {
  mcp?: string;
  tool?: string;
  connectionUuid?: string;
  params?: Record<string, string>;
};

export type ScriptStep = {
  language?: ScriptLanguage;
  code?: string;
  inputParams?: string;
  outputParams?: string;
  params?: Record<string, string>;
};

export type ConditionStep = {
  conditions?: TractCondition[];
  then?: TractStep[];
  else?: TractStep[];
};

export type ParallelStep = {
  steps?: TractStep[];
};

export type GroupStep = {
  steps?: TractStep[];
};

type BaseTractStep = {
  id?: string;
  name?: string;
  description?: string;
};

export type TractStep = BaseTractStep &
  OneOf<{
    action: ActionStep;
    condition: ConditionStep;
    parallel: ParallelStep;
    group: GroupStep;
    script: ScriptStep;
  }>;

export type TractDefinition = {
  steps?: TractStep[];
};

export type TractItem = {
  uuid?: string;
  name?: string;
  description?: string;
  enabled?: boolean;
  definition?: TractDefinition;
  triggers?: TractTriggerSummary[];
  lastRun?: TractLastRun;
  createdAt?: string;
  updatedAt?: string;
};

export type TractRunItem = {
  uuid?: string;
  tractUuid?: string;
  triggerUuid?: string;
  status?: string;
  startedBy?: string;
  triggerPayload?: string;
  error?: string;
  createdAt?: string;
  updatedAt?: string;
};

export type TractRunStepItem = {
  stepId?: string;
  stepName?: string;
  stepType?: string;
  status?: string;
  input?: string;
  output?: string;
  error?: string;
  startedAt?: string;
  finishedAt?: string;
};

export type TractToolItem = {
  mcp?: string;
  tool?: string;
  description?: string;
  inputSchema?: string;
  outputSchema?: string;
};

export type TriggerSourceItem = {
  key?: string;
  description?: string;
  payloadSchema?: string;
  category?: string;
  label?: string;
  provider?: string;
};

export type TriggerItem = {
  uuid?: string;
  name?: string;
  kind?: string;
  source?: string;
  config?: string;
  payloadSchema?: string;
  triggerUuid?: string;
  enabled?: boolean;
  createdAt?: string;
  tokenSuffix?: string;
};

export type CreateTractRequest = {
  name?: string;
  description?: string;
  definition?: TractDefinition;
};

export type CreateTractResponse = {
  tract?: TractItem;
  warnings?: string[];
};

export type CreateTract = Record<string, never>;

export type UpdateTractRequest = {
  uuid?: string;
  name?: string;
  description?: string;
  definition?: TractDefinition;
};

export type UpdateTractResponse = {
  tract?: TractItem;
  warnings?: string[];
};

export type UpdateTract = Record<string, never>;

export type GetTractRequest = {
  uuid?: string;
};

export type GetTractResponse = {
  tract?: TractItem;
};

export type GetTract = Record<string, never>;

export type ListTractsRequest = Record<string, never>;

export type ListTractsResponse = {
  tracts?: TractItem[];
};

export type ListTracts = Record<string, never>;

export type DeleteTractRequest = {
  uuid?: string;
};

export type DeleteTractResponse = Record<string, never>;

export type DeleteTract = Record<string, never>;

export type SetTractEnabledRequest = {
  uuid?: string;
  enabled?: boolean;
};

export type SetTractEnabledResponse = Record<string, never>;

export type SetTractEnabled = Record<string, never>;

export type RunTractRequest = {
  tractUuid?: string;
  params?: string;
};

export type RunTractResponse = Record<string, never>;

export type RunTract = Record<string, never>;

export type ListRunsRequest = {
  tractUuid?: string;
  limit?: number;
};

export type ListRunsResponse = {
  runs?: TractRunItem[];
};

export type ListRuns = Record<string, never>;

export type GetRunRequest = {
  runUuid?: string;
};

export type GetRunResponse = {
  run?: TractRunItem;
  steps?: TractRunStepItem[];
};

export type GetRun = Record<string, never>;

export type WatchRunRequest = {
  runUuid?: string;
};

export type WatchRunResponse = {
  run?: TractRunItem;
  steps?: TractRunStepItem[];
};

export type WatchRun = Record<string, never>;

export type RetryRunRequest = {
  runUuid?: string;
};

export type RetryRunResponse = Record<string, never>;

export type RetryRun = Record<string, never>;

export type ListTractToolsRequest = Record<string, never>;

export type ListTractToolsResponse = {
  tools?: TractToolItem[];
};

export type ListTractTools = Record<string, never>;

export type ListTriggerSourcesRequest = Record<string, never>;

export type ListTriggerSourcesResponse = {
  sources?: TriggerSourceItem[];
};

export type ListTriggerSources = Record<string, never>;

export type CreateTriggerRequest = {
  name?: string;
  kind?: string;
  source?: string;
  config?: string;
  payloadSchema?: string;
};

export type CreateTriggerResponse = {
  trigger?: TriggerItem;
  webhookUrl?: string;
  webhookToken?: string;
};

export type CreateTrigger = Record<string, never>;

export type ListTriggersRequest = Record<string, never>;

export type ListTriggersResponse = {
  triggers?: TriggerItem[];
};

export type ListTriggers = Record<string, never>;

export type DeleteTriggerRequest = {
  uuid?: string;
};

export type DeleteTriggerResponse = Record<string, never>;

export type DeleteTrigger = Record<string, never>;

export type SetTriggerEnabledRequest = {
  uuid?: string;
  enabled?: boolean;
};

export type SetTriggerEnabledResponse = Record<string, never>;

export type SetTriggerEnabled = Record<string, never>;

export type RotateTriggerTokenRequest = {
  uuid?: string;
};

export type RotateTriggerTokenResponse = {
  trigger?: TriggerItem;
  webhookUrl?: string;
  webhookToken?: string;
};

export type RotateTriggerToken = Record<string, never>;

export type LinkTriggerRequest = {
  triggerUuid?: string;
  tractUuid?: string;
  filters?: string;
};

export type LinkTriggerResponse = Record<string, never>;

export type LinkTrigger = Record<string, never>;

export type UnlinkTriggerRequest = {
  triggerUuid?: string;
  tractUuid?: string;
};

export type UnlinkTriggerResponse = Record<string, never>;

export type UnlinkTrigger = Record<string, never>;

export class TractsAPI {
  static CreateTract(this:void, req: CreateTractRequest, initReq?: fm.InitReq): Promise<CreateTractResponse> {
    return fm.fetchRequest<CreateTractResponse>(`/api/tracts/create`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static UpdateTract(this:void, req: UpdateTractRequest, initReq?: fm.InitReq): Promise<UpdateTractResponse> {
    return fm.fetchRequest<UpdateTractResponse>(`/api/tracts/update`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static GetTract(this:void, req: GetTractRequest, initReq?: fm.InitReq): Promise<GetTractResponse> {
    return fm.fetchRequest<GetTractResponse>(`/api/tracts/get`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ListTracts(this:void, req: ListTractsRequest, initReq?: fm.InitReq): Promise<ListTractsResponse> {
    return fm.fetchRequest<ListTractsResponse>(`/api/tracts/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static DeleteTract(this:void, req: DeleteTractRequest, initReq?: fm.InitReq): Promise<DeleteTractResponse> {
    return fm.fetchRequest<DeleteTractResponse>(`/api/tracts/delete`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static SetTractEnabled(this:void, req: SetTractEnabledRequest, initReq?: fm.InitReq): Promise<SetTractEnabledResponse> {
    return fm.fetchRequest<SetTractEnabledResponse>(`/api/tracts/set-enabled`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static RunTract(this:void, req: RunTractRequest, initReq?: fm.InitReq): Promise<RunTractResponse> {
    return fm.fetchRequest<RunTractResponse>(`/api/tracts/run`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ListRuns(this:void, req: ListRunsRequest, initReq?: fm.InitReq): Promise<ListRunsResponse> {
    return fm.fetchRequest<ListRunsResponse>(`/api/tracts/runs/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static GetRun(this:void, req: GetRunRequest, initReq?: fm.InitReq): Promise<GetRunResponse> {
    return fm.fetchRequest<GetRunResponse>(`/api/tracts/runs/get`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static WatchRun(this:void, req: WatchRunRequest, entityNotifier?: fm.NotifyStreamEntityArrival<WatchRunResponse>, initReq?: fm.InitReq): Promise<void> {
    return fm.fetchStreamingRequest<WatchRunResponse>(`/api/tracts/runs/watch`, entityNotifier, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static RetryRun(this:void, req: RetryRunRequest, initReq?: fm.InitReq): Promise<RetryRunResponse> {
    return fm.fetchRequest<RetryRunResponse>(`/api/tracts/runs/retry`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ListTractTools(this:void, req: ListTractToolsRequest, initReq?: fm.InitReq): Promise<ListTractToolsResponse> {
    return fm.fetchRequest<ListTractToolsResponse>(`/api/tracts/tools/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ListTriggerSources(this:void, req: ListTriggerSourcesRequest, initReq?: fm.InitReq): Promise<ListTriggerSourcesResponse> {
    return fm.fetchRequest<ListTriggerSourcesResponse>(`/api/tracts/trigger-sources/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static CreateTrigger(this:void, req: CreateTriggerRequest, initReq?: fm.InitReq): Promise<CreateTriggerResponse> {
    return fm.fetchRequest<CreateTriggerResponse>(`/api/tracts/triggers/create`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ListTriggers(this:void, req: ListTriggersRequest, initReq?: fm.InitReq): Promise<ListTriggersResponse> {
    return fm.fetchRequest<ListTriggersResponse>(`/api/tracts/triggers/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static DeleteTrigger(this:void, req: DeleteTriggerRequest, initReq?: fm.InitReq): Promise<DeleteTriggerResponse> {
    return fm.fetchRequest<DeleteTriggerResponse>(`/api/tracts/triggers/delete`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static SetTriggerEnabled(this:void, req: SetTriggerEnabledRequest, initReq?: fm.InitReq): Promise<SetTriggerEnabledResponse> {
    return fm.fetchRequest<SetTriggerEnabledResponse>(`/api/tracts/triggers/set-enabled`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static RotateTriggerToken(this:void, req: RotateTriggerTokenRequest, initReq?: fm.InitReq): Promise<RotateTriggerTokenResponse> {
    return fm.fetchRequest<RotateTriggerTokenResponse>(`/api/tracts/triggers/${req.uuid}/rotate_token`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static LinkTrigger(this:void, req: LinkTriggerRequest, initReq?: fm.InitReq): Promise<LinkTriggerResponse> {
    return fm.fetchRequest<LinkTriggerResponse>(`/api/tracts/triggers/link`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static UnlinkTrigger(this:void, req: UnlinkTriggerRequest, initReq?: fm.InitReq): Promise<UnlinkTriggerResponse> {
    return fm.fetchRequest<UnlinkTriggerResponse>(`/api/tracts/triggers/unlink`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
}