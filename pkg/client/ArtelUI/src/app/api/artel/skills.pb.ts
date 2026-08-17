/* eslint-disable */
// @ts-nocheck

/**
 * This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
 */

import * as fm from "./fetch.pb";


export type SkillInfo = {
  slug?: string;
  name?: string;
  description?: string;
  storageMode?: string;
  isHotPlug?: boolean;
  isSystem?: boolean;
};

export type ListSkillsRequest = {
  vaultId?: string;
};

export type ListSkillsResponse = {
  skills?: SkillInfo[];
};

export type ListSkills = Record<string, never>;

export type GetSkillRequest = {
  vaultId?: string;
  slug?: string;
};

export type GetSkillResponse = {
  skill?: SkillInfo;
  body?: string;
};

export type GetSkill = Record<string, never>;

export type CreateSkillRequest = {
  vaultId?: string;
  name?: string;
  description?: string;
  storageMode?: string;
  body?: string;
  hotPlug?: boolean;
};

export type CreateSkillResponse = {
  skill?: SkillInfo;
};

export type CreateSkill = Record<string, never>;

export type UpdateSkillRequest = {
  vaultId?: string;
  slug?: string;
  name?: string;
  description?: string;
  storageMode?: string;
  body?: string;
};

export type UpdateSkillResponse = {
  skill?: SkillInfo;
};

export type UpdateSkill = Record<string, never>;

export type DeleteSkillRequest = {
  vaultId?: string;
  slug?: string;
};

export type DeleteSkillResponse = Record<string, never>;

export type DeleteSkill = Record<string, never>;

export type SetSkillHotPlugRequest = {
  vaultId?: string;
  slug?: string;
  hotPlug?: boolean;
};

export type SetSkillHotPlugResponse = {
  skill?: SkillInfo;
};

export type SetSkillHotPlug = Record<string, never>;

export class SkillsAPI {
  static ListSkills(this:void, req: ListSkillsRequest, initReq?: fm.InitReq): Promise<ListSkillsResponse> {
    return fm.fetchRequest<ListSkillsResponse>(`/api/skills/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static GetSkill(this:void, req: GetSkillRequest, initReq?: fm.InitReq): Promise<GetSkillResponse> {
    return fm.fetchRequest<GetSkillResponse>(`/api/skills/get`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static CreateSkill(this:void, req: CreateSkillRequest, initReq?: fm.InitReq): Promise<CreateSkillResponse> {
    return fm.fetchRequest<CreateSkillResponse>(`/api/skills/create`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static UpdateSkill(this:void, req: UpdateSkillRequest, initReq?: fm.InitReq): Promise<UpdateSkillResponse> {
    return fm.fetchRequest<UpdateSkillResponse>(`/api/skills/update`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static DeleteSkill(this:void, req: DeleteSkillRequest, initReq?: fm.InitReq): Promise<DeleteSkillResponse> {
    return fm.fetchRequest<DeleteSkillResponse>(`/api/skills/delete`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static SetSkillHotPlug(this:void, req: SetSkillHotPlugRequest, initReq?: fm.InitReq): Promise<SetSkillHotPlugResponse> {
    return fm.fetchRequest<SetSkillHotPlugResponse>(`/api/skills/hot-plug`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
}