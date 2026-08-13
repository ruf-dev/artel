/* eslint-disable */
// @ts-nocheck

/**
 * This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
 */

import * as fm from "./fetch.pb";
import * as ArtelPublicDocsPublicDocs from "./public_docs.pb";
import * as ArtelSetupWizardSetupWizard from "./setup_wizard.pb";


export type GetSettingsRequest = Record<string, never>;

export type GetSettingsResponse = {
  passwordAuthEnabled?: boolean;
  telegramAuthEnabled?: boolean;
  registrationMode?: ArtelSetupWizardSetupWizard.RegistrationMode;
  defaultDocsVaultId?: string;
  defaultDocsVaultName?: string;
  defaultDocsVaultSlug?: string;
  defaultDocsSource?: ArtelPublicDocsPublicDocs.DocsSource;
};

export type GetSettings = Record<string, never>;

export type UpdateAuthMethodsRequest = {
  passwordEnabled?: boolean;
  telegramEnabled?: boolean;
};

export type UpdateAuthMethodsResponse = Record<string, never>;

export type UpdateAuthMethods = Record<string, never>;

export type UpdateRegistrationModeRequest = {
  mode?: ArtelSetupWizardSetupWizard.RegistrationMode;
};

export type UpdateRegistrationModeResponse = Record<string, never>;

export type UpdateRegistrationMode = Record<string, never>;

export type UpdateDefaultDocsVaultRequest = {
  vaultId?: string;
  source?: ArtelPublicDocsPublicDocs.DocsSource;
};

export type UpdateDefaultDocsVaultResponse = Record<string, never>;

export type UpdateDefaultDocsVault = Record<string, never>;

export class AdminSystemSettingsAPI {
  static GetSettings(this:void, req: GetSettingsRequest, initReq?: fm.InitReq): Promise<GetSettingsResponse> {
    return fm.fetchRequest<GetSettingsResponse>(`/api/admin_system_settings/get`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static UpdateAuthMethods(this:void, req: UpdateAuthMethodsRequest, initReq?: fm.InitReq): Promise<UpdateAuthMethodsResponse> {
    return fm.fetchRequest<UpdateAuthMethodsResponse>(`/api/admin_system_settings/update_auth_methods`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static UpdateRegistrationMode(this:void, req: UpdateRegistrationModeRequest, initReq?: fm.InitReq): Promise<UpdateRegistrationModeResponse> {
    return fm.fetchRequest<UpdateRegistrationModeResponse>(`/api/admin_system_settings/update_registration_mode`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static UpdateDefaultDocsVault(this:void, req: UpdateDefaultDocsVaultRequest, initReq?: fm.InitReq): Promise<UpdateDefaultDocsVaultResponse> {
    return fm.fetchRequest<UpdateDefaultDocsVaultResponse>(`/api/admin_system_settings/update_default_docs_vault`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
}