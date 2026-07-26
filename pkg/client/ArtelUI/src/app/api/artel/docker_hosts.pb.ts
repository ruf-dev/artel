/* eslint-disable */
// @ts-nocheck

/**
 * This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
 */

import * as fm from "./fetch.pb";


export type RegisterDockerHostRequest = {
  url?: string;
  caCert?: string;
  clientCert?: string;
  clientKey?: string;
};

export type RegisterDockerHostResponse = {
  id?: string;
};

export type RegisterDockerHost = Record<string, never>;

export type GetDockerHostRequest = {
  id?: string;
};

export type GetDockerHostResponse = {
  id?: string;
  url?: string;
  createdAt?: string;
};

export type GetDockerHost = Record<string, never>;

export type ListDockerHostsRequest = Record<string, never>;

export type ListDockerHostsResponse = {
  hosts?: GetDockerHostResponse[];
};

export type ListDockerHosts = Record<string, never>;

export type UpdateDockerHostRequest = {
  id?: string;
  url?: string;
  caCert?: string;
  clientCert?: string;
  clientKey?: string;
};

export type UpdateDockerHostResponse = Record<string, never>;

export type UpdateDockerHost = Record<string, never>;

export type DeleteDockerHostRequest = {
  id?: string;
};

export type DeleteDockerHostResponse = Record<string, never>;

export type DeleteDockerHost = Record<string, never>;

export class DockerHostsAPI {
  static RegisterDockerHost(this:void, req: RegisterDockerHostRequest, initReq?: fm.InitReq): Promise<RegisterDockerHostResponse> {
    return fm.fetchRequest<RegisterDockerHostResponse>(`/api/docker_hosts/add`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static GetDockerHost(this:void, req: GetDockerHostRequest, initReq?: fm.InitReq): Promise<GetDockerHostResponse> {
    return fm.fetchRequest<GetDockerHostResponse>(`/api/docker_hosts/get`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ListDockerHosts(this:void, req: ListDockerHostsRequest, initReq?: fm.InitReq): Promise<ListDockerHostsResponse> {
    return fm.fetchRequest<ListDockerHostsResponse>(`/api/docker_hosts/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static UpdateDockerHost(this:void, req: UpdateDockerHostRequest, initReq?: fm.InitReq): Promise<UpdateDockerHostResponse> {
    return fm.fetchRequest<UpdateDockerHostResponse>(`/api/docker_hosts/update`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static DeleteDockerHost(this:void, req: DeleteDockerHostRequest, initReq?: fm.InitReq): Promise<DeleteDockerHostResponse> {
    return fm.fetchRequest<DeleteDockerHostResponse>(`/api/docker_hosts/delete`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
}