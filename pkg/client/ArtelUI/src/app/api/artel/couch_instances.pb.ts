/* eslint-disable */
// @ts-nocheck

/**
 * This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
 */

import * as fm from "./fetch.pb";


export type RegisterCouchInstanceRequest = {
  url?: string;
  username?: string;
  password?: string;
};

export type RegisterCouchInstanceResponse = {
  id?: string;
};

export type RegisterCouchInstance = Record<string, never>;

export type GetCouchInstanceRequest = {
  id?: string;
};

export type GetCouchInstanceResponse = {
  id?: string;
  url?: string;
  username?: string;
  createdAt?: string;
};

export type GetCouchInstance = Record<string, never>;

export type ListCouchInstancesRequest = Record<string, never>;

export type ListCouchInstancesResponse = {
  instances?: GetCouchInstanceResponse[];
};

export type ListCouchInstances = Record<string, never>;

export type UpdateCouchInstanceRequest = {
  id?: string;
  url?: string;
  username?: string;
  password?: string;
};

export type UpdateCouchInstanceResponse = Record<string, never>;

export type UpdateCouchInstance = Record<string, never>;

export type DeleteCouchInstanceRequest = {
  id?: string;
};

export type DeleteCouchInstanceResponse = Record<string, never>;

export type DeleteCouchInstance = Record<string, never>;

export type SetupCouchInstanceRequest = {
  id?: string;
};

export type SetupCouchInstanceResponse = Record<string, never>;

export type SetupCouchInstance = Record<string, never>;

export type GetCouchInstanceStatusRequest = {
  id?: string;
};

export type GetCouchInstanceStatusResponse = {
  clusterModeEnabled?: boolean;
  usersDbInitialized?: boolean;
  replicatorDbInitialized?: boolean;
};

export type GetCouchInstanceStatus = Record<string, never>;

export class CouchInstancesAPI {
  static RegisterCouchInstance(this:void, req: RegisterCouchInstanceRequest, initReq?: fm.InitReq): Promise<RegisterCouchInstanceResponse> {
    return fm.fetchRequest<RegisterCouchInstanceResponse>(`/api/couch/add`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static GetCouchInstance(this:void, req: GetCouchInstanceRequest, initReq?: fm.InitReq): Promise<GetCouchInstanceResponse> {
    return fm.fetchRequest<GetCouchInstanceResponse>(`/api/couch/get`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ListCouchInstances(this:void, req: ListCouchInstancesRequest, initReq?: fm.InitReq): Promise<ListCouchInstancesResponse> {
    return fm.fetchRequest<ListCouchInstancesResponse>(`/api/couch/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static UpdateCouchInstance(this:void, req: UpdateCouchInstanceRequest, initReq?: fm.InitReq): Promise<UpdateCouchInstanceResponse> {
    return fm.fetchRequest<UpdateCouchInstanceResponse>(`/api/couch/update`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static DeleteCouchInstance(this:void, req: DeleteCouchInstanceRequest, initReq?: fm.InitReq): Promise<DeleteCouchInstanceResponse> {
    return fm.fetchRequest<DeleteCouchInstanceResponse>(`/api/couch/delete`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static SetupCouchInstance(this:void, req: SetupCouchInstanceRequest, initReq?: fm.InitReq): Promise<SetupCouchInstanceResponse> {
    return fm.fetchRequest<SetupCouchInstanceResponse>(`/api/couch/setup`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static GetCouchInstanceStatus(this:void, req: GetCouchInstanceStatusRequest, initReq?: fm.InitReq): Promise<GetCouchInstanceStatusResponse> {
    return fm.fetchRequest<GetCouchInstanceStatusResponse>(`/api/couch/status`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
}