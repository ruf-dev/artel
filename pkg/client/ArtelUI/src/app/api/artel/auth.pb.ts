/* eslint-disable */
// @ts-nocheck

/**
 * This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
 */

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

export type PasswordCredentials = {
  email?: string;
  password?: string;
};

export type TelegramCredentials = {
  idToken?: string;
};

export type RegisterRequest = {
  email?: string;
  password?: string;
};

export type RegisterResponse = {
  id?: string;
  email?: string;
};

export type Register = Record<string, never>;

type BaseLoginRequest = {
};

export type LoginRequest = BaseLoginRequest &
  OneOf<{
    password: PasswordCredentials;
    telegram: TelegramCredentials;
  }>;

export type LoginResponse = {
  token?: string;
  expiresAt?: GoogleProtobufTimestamp.Timestamp;
};

export type Login = Record<string, never>;

export type LogoutRequest = Record<string, never>;

export type LogoutResponse = Record<string, never>;

export type Logout = Record<string, never>;

export type GetConfigRequest = Record<string, never>;

export type GetConfigResponse = {
  telegramClientId?: string;
};

export type GetConfig = Record<string, never>;

export type Permissions = {
  isAdministrator?: boolean;
  hasEmails?: boolean;
};

export type GetMeRequest = Record<string, never>;

export type GetMeResponse = {
  id?: string;
  username?: string;
  email?: string;
  permissions?: Permissions;
  photoUrl?: string;
};

export type GetMe = Record<string, never>;

export class AuthAPI {
  static Register(this:void, req: RegisterRequest, initReq?: fm.InitReq): Promise<RegisterResponse> {
    return fm.fetchRequest<RegisterResponse>(`/api/auth/register`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static Login(this:void, req: LoginRequest, initReq?: fm.InitReq): Promise<LoginResponse> {
    return fm.fetchRequest<LoginResponse>(`/api/auth/login`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static Logout(this:void, req: LogoutRequest, initReq?: fm.InitReq): Promise<LogoutResponse> {
    return fm.fetchRequest<LogoutResponse>(`/api/auth/logout`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static GetConfig(this:void, req: GetConfigRequest, initReq?: fm.InitReq): Promise<GetConfigResponse> {
    return fm.fetchRequest<GetConfigResponse>(`/api/auth/config?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"});
  }
  static GetMe(this:void, req: GetMeRequest, initReq?: fm.InitReq): Promise<GetMeResponse> {
    return fm.fetchRequest<GetMeResponse>(`/api/auth/me?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"});
  }
}