/* eslint-disable */
// @ts-nocheck

/**
 * This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
 */

import * as fm from "./fetch.pb";


export enum PromptId {
  PROMPT_ID_UNSPECIFIED = "PROMPT_ID_UNSPECIFIED",
  PROMPT_ID_PERSONAL_VAULT_SETUP = "PROMPT_ID_PERSONAL_VAULT_SETUP",
}

export type PromptItem = {
  id?: PromptId;
  prompt?: string;
};

export type ListPromptsRequest = {
  ids?: PromptId[];
  page?: number;
  pageSize?: number;
};

export type ListPromptsResponse = {
  prompts?: PromptItem[];
  total?: number;
};

export type ListPrompts = Record<string, never>;

export class PromptsAPI {
  static ListPrompts(this:void, req: ListPromptsRequest, initReq?: fm.InitReq): Promise<ListPromptsResponse> {
    return fm.fetchRequest<ListPromptsResponse>(`/api/prompts/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
}