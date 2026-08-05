/* eslint-disable */
// @ts-nocheck

/**
 * This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
 */

import * as ArtelAuthAuth from "./auth.pb";
import * as fm from "./fetch.pb";
import * as GoogleProtobufTimestamp from "./google/protobuf/timestamp.pb";

type Absent<T, K extends keyof T> = { [k in Exclude<keyof T, K>]?: undefined };

type OneOf<T> =
  | { [k in keyof T]?: undefined }
  | (keyof T extends infer K
      ? K extends string & keyof T
        ? { [k in K]: T[K] } & Absent<T, K>
        : never
      : never);

export enum RegistrationMode {
  ADMIN_ONLY = "ADMIN_ONLY",
  SELF_REGISTER = "SELF_REGISTER",
}

export type GetStatusRequest = Record<string, never>;

export type GetStatusResponse = {
  setupCompleted?: boolean;
  tokenPending?: boolean;
};

export type GetStatus = Record<string, never>;

export type SubmitTokenRequest = {
  token?: string;
};

export type SubmitTokenResponse = {
  wizardSessionToken?: string;
};

export type SubmitToken = Record<string, never>;

export type SelectAuthMethodsRequest = {
  wizardSessionToken?: string;
  passwordEnabled?: boolean;
  telegramEnabled?: boolean;
};

export type SelectAuthMethodsResponse = Record<string, never>;

export type SelectAuthMethods = Record<string, never>;

export type SelectRegistrationModeRequest = {
  wizardSessionToken?: string;
  mode?: RegistrationMode;
};

export type SelectRegistrationModeResponse = Record<string, never>;

export type SelectRegistrationMode = Record<string, never>;

type BaseCompleteSetupRequest = {
  wizardSessionToken?: string;
};

export type CompleteSetupRequest = BaseCompleteSetupRequest &
  OneOf<{
    password: ArtelAuthAuth.PasswordCredentials;
    telegram: ArtelAuthAuth.TelegramCredentials;
  }>;

export type CompleteSetupResponse = {
  token?: string;
  expiresAt?: GoogleProtobufTimestamp.Timestamp;
  refreshToken?: string;
  refreshExpiresAt?: GoogleProtobufTimestamp.Timestamp;
};

export type CompleteSetup = Record<string, never>;

export class SetupWizardAPI {
  static GetStatus(this:void, req: GetStatusRequest, initReq?: fm.InitReq): Promise<GetStatusResponse> {
    return fm.fetchRequest<GetStatusResponse>(`/api/setup_wizard/status`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static SubmitToken(this:void, req: SubmitTokenRequest, initReq?: fm.InitReq): Promise<SubmitTokenResponse> {
    return fm.fetchRequest<SubmitTokenResponse>(`/api/setup_wizard/submit_token`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static SelectAuthMethods(this:void, req: SelectAuthMethodsRequest, initReq?: fm.InitReq): Promise<SelectAuthMethodsResponse> {
    return fm.fetchRequest<SelectAuthMethodsResponse>(`/api/setup_wizard/select_auth_methods`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static SelectRegistrationMode(this:void, req: SelectRegistrationModeRequest, initReq?: fm.InitReq): Promise<SelectRegistrationModeResponse> {
    return fm.fetchRequest<SelectRegistrationModeResponse>(`/api/setup_wizard/select_registration_mode`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static CompleteSetup(this:void, req: CompleteSetupRequest, initReq?: fm.InitReq): Promise<CompleteSetupResponse> {
    return fm.fetchRequest<CompleteSetupResponse>(`/api/setup_wizard/complete`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
}