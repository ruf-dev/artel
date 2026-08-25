/* eslint-disable */
// @ts-nocheck

/**
 * This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
 */

import * as fm from "./fetch.pb";


export type SimpleChat = {
  id?: string;
  vaultId?: string;
  title?: string;
  model?: string;
  vaultAccess?: boolean;
  createdAt?: string;
  updatedAt?: string;
  lastActivityAt?: string;
};

export type SimpleChatMessage = {
  id?: string;
  role?: string;
  content?: string;
  toolName?: string;
  isError?: boolean;
  model?: string;
  seq?: string;
  createdAt?: string;
};

export type CreateSimpleChatRequest = {
  vaultId?: string;
  vaultAccess?: boolean;
  model?: string;
};

export type CreateSimpleChatResponse = {
  chat?: SimpleChat;
};

export type CreateSimpleChat = Record<string, never>;

export type ListSimpleChatsRequest = {
  vaultId?: string;
};

export type ListSimpleChatsResponse = {
  chats?: SimpleChat[];
};

export type ListSimpleChats = Record<string, never>;

export type GetSimpleChatRequest = {
  chatId?: string;
};

export type GetSimpleChatResponse = {
  chat?: SimpleChat;
  messages?: SimpleChatMessage[];
};

export type GetSimpleChat = Record<string, never>;

export type DeleteSimpleChatRequest = {
  chatId?: string;
};

export type DeleteSimpleChatResponse = Record<string, never>;

export type DeleteSimpleChat = Record<string, never>;

export class SimpleChatAPI {
  static CreateSimpleChat(this:void, req: CreateSimpleChatRequest, initReq?: fm.InitReq): Promise<CreateSimpleChatResponse> {
    return fm.fetchRequest<CreateSimpleChatResponse>(`/api/simple-chats/create`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ListSimpleChats(this:void, req: ListSimpleChatsRequest, initReq?: fm.InitReq): Promise<ListSimpleChatsResponse> {
    return fm.fetchRequest<ListSimpleChatsResponse>(`/api/simple-chats/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static GetSimpleChat(this:void, req: GetSimpleChatRequest, initReq?: fm.InitReq): Promise<GetSimpleChatResponse> {
    return fm.fetchRequest<GetSimpleChatResponse>(`/api/simple-chats/get`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static DeleteSimpleChat(this:void, req: DeleteSimpleChatRequest, initReq?: fm.InitReq): Promise<DeleteSimpleChatResponse> {
    return fm.fetchRequest<DeleteSimpleChatResponse>(`/api/simple-chats/delete`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
}