/* eslint-disable */
// @ts-nocheck

/**
 * This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
 */

import * as fm from "./fetch.pb";


export type GetUserSettingsRequest = Record<string, never>;

export type GetUserSettingsResponse = {
  userPrompt?: string;
  likedOpenrouterModels?: string[];
  lastUsedModel?: string;
};

export type GetUserSettings = Record<string, never>;

export type SetLikedModelsRequest = {
  likedOpenrouterModels?: string[];
};

export type SetLikedModelsResponse = Record<string, never>;

export type SetLikedModels = Record<string, never>;

export type SetLastUsedModelRequest = {
  model?: string;
};

export type SetLastUsedModelResponse = Record<string, never>;

export type SetLastUsedModel = Record<string, never>;

export class UserSettingsAPI {
  static GetUserSettings(this:void, req: GetUserSettingsRequest, initReq?: fm.InitReq): Promise<GetUserSettingsResponse> {
    return fm.fetchRequest<GetUserSettingsResponse>(`/api/user-settings/get?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"});
  }
  static SetLikedModels(this:void, req: SetLikedModelsRequest, initReq?: fm.InitReq): Promise<SetLikedModelsResponse> {
    return fm.fetchRequest<SetLikedModelsResponse>(`/api/user-settings/set-liked-models`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static SetLastUsedModel(this:void, req: SetLastUsedModelRequest, initReq?: fm.InitReq): Promise<SetLastUsedModelResponse> {
    return fm.fetchRequest<SetLastUsedModelResponse>(`/api/user-settings/set-last-used-model`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
}