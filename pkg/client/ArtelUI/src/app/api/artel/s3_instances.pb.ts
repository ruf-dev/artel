/* eslint-disable */
// @ts-nocheck

/**
 * This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
 */

import * as fm from "./fetch.pb";


export type RegisterS3InstanceRequest = {
  endpoint?: string;
  region?: string;
  accessKey?: string;
  secretKey?: string;
  useSsl?: boolean;
  pathStyle?: boolean;
};

export type RegisterS3InstanceResponse = {
  id?: string;
};

export type RegisterS3Instance = Record<string, never>;

export type GetS3InstanceRequest = {
  id?: string;
};

export type GetS3InstanceResponse = {
  id?: string;
  endpoint?: string;
  region?: string;
  useSsl?: boolean;
  pathStyle?: boolean;
  createdAt?: string;
};

export type GetS3Instance = Record<string, never>;

export type ListS3InstancesRequest = Record<string, never>;

export type ListS3InstancesResponse = {
  instances?: GetS3InstanceResponse[];
};

export type ListS3Instances = Record<string, never>;

export type UpdateS3InstanceRequest = {
  id?: string;
  endpoint?: string;
  region?: string;
  accessKey?: string;
  secretKey?: string;
  useSsl?: boolean;
  pathStyle?: boolean;
};

export type UpdateS3InstanceResponse = Record<string, never>;

export type UpdateS3Instance = Record<string, never>;

export type DeleteS3InstanceRequest = {
  id?: string;
};

export type DeleteS3InstanceResponse = Record<string, never>;

export type DeleteS3Instance = Record<string, never>;

export type TestS3InstanceRequest = {
  id?: string;
};

export type TestS3InstanceResponse = Record<string, never>;

export type TestS3Instance = Record<string, never>;

export class S3InstancesAPI {
  static RegisterS3Instance(this:void, req: RegisterS3InstanceRequest, initReq?: fm.InitReq): Promise<RegisterS3InstanceResponse> {
    return fm.fetchRequest<RegisterS3InstanceResponse>(`/api/s3/add`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static GetS3Instance(this:void, req: GetS3InstanceRequest, initReq?: fm.InitReq): Promise<GetS3InstanceResponse> {
    return fm.fetchRequest<GetS3InstanceResponse>(`/api/s3/get`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ListS3Instances(this:void, req: ListS3InstancesRequest, initReq?: fm.InitReq): Promise<ListS3InstancesResponse> {
    return fm.fetchRequest<ListS3InstancesResponse>(`/api/s3/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static UpdateS3Instance(this:void, req: UpdateS3InstanceRequest, initReq?: fm.InitReq): Promise<UpdateS3InstanceResponse> {
    return fm.fetchRequest<UpdateS3InstanceResponse>(`/api/s3/update`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static DeleteS3Instance(this:void, req: DeleteS3InstanceRequest, initReq?: fm.InitReq): Promise<DeleteS3InstanceResponse> {
    return fm.fetchRequest<DeleteS3InstanceResponse>(`/api/s3/delete`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static TestS3Instance(this:void, req: TestS3InstanceRequest, initReq?: fm.InitReq): Promise<TestS3InstanceResponse> {
    return fm.fetchRequest<TestS3InstanceResponse>(`/api/s3/test`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
}