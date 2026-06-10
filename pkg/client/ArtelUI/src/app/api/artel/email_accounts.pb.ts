/* eslint-disable */
// @ts-nocheck

/**
 * This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
 */

import * as fm from "./fetch.pb";


export type EmailAccountInfo = {
  id?: string;
  email?: string;
  imapHost?: string;
  imapPort?: number;
  smtpHost?: string;
  smtpPort?: number;
  createdAt?: string;
};

export type AddEmailAccountRequest = {
  email?: string;
  imapHost?: string;
  imapPort?: number;
  smtpHost?: string;
  smtpPort?: number;
  password?: string;
};

export type AddEmailAccountResponse = {
  account?: EmailAccountInfo;
};

export type AddEmailAccount = Record<string, never>;

export type ListEmailAccountsRequest = Record<string, never>;

export type ListEmailAccountsResponse = {
  accounts?: EmailAccountInfo[];
};

export type ListEmailAccounts = Record<string, never>;

export type DeleteEmailAccountRequest = {
  id?: string;
};

export type DeleteEmailAccountResponse = Record<string, never>;

export type DeleteEmailAccount = Record<string, never>;

export type MailServerSuggestionInfo = {
  domain?: string;
  smtp?: string;
  smtpPort?: number;
  imap?: string;
  imapPort?: number;
};

export type ListMailServerSuggestionsRequest = {
  domain?: string;
};

export type ListMailServerSuggestionsResponse = {
  suggestions?: MailServerSuggestionInfo[];
};

export type ListMailServerSuggestions = Record<string, never>;

export class EmailAccountsAPI {
  static AddEmailAccount(this:void, req: AddEmailAccountRequest, initReq?: fm.InitReq): Promise<AddEmailAccountResponse> {
    return fm.fetchRequest<AddEmailAccountResponse>(`/api/email-accounts/add`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ListEmailAccounts(this:void, req: ListEmailAccountsRequest, initReq?: fm.InitReq): Promise<ListEmailAccountsResponse> {
    return fm.fetchRequest<ListEmailAccountsResponse>(`/api/email-accounts/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static DeleteEmailAccount(this:void, req: DeleteEmailAccountRequest, initReq?: fm.InitReq): Promise<DeleteEmailAccountResponse> {
    return fm.fetchRequest<DeleteEmailAccountResponse>(`/api/email-accounts/delete`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ListMailServerSuggestions(this:void, req: ListMailServerSuggestionsRequest, initReq?: fm.InitReq): Promise<ListMailServerSuggestionsResponse> {
    return fm.fetchRequest<ListMailServerSuggestionsResponse>(`/api/email-accounts/list_suggestions`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
}