/* eslint-disable */
// @ts-nocheck

/**
 * This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
 */

import * as fm from "./fetch.pb";


export type SubscriptionPlanEntry = {
  planKey?: string;
  couchQuotaBytes?: string;
  s3QuotaBytes?: string;
  features?: Record<string, boolean>;
};

export type ListSubscriptionPlansRequest = Record<string, never>;

export type ListSubscriptionPlansResponse = {
  plans?: SubscriptionPlanEntry[];
};

export type ListSubscriptionPlans = Record<string, never>;

export type EffectiveSubscriptionView = {
  active?: boolean;
  planKey?: string;
  features?: Record<string, boolean>;
  couchQuotaBytes?: string;
  s3QuotaBytes?: string;
};

export type GetUserSubscriptionRequest = {
  userId?: string;
};

export type GetUserSubscriptionResponse = {
  active?: boolean;
  planKey?: string;
  featureOverrides?: Record<string, boolean>;
  couchQuotaOverrideBytes?: string;
  s3QuotaOverrideBytes?: string;
  effective?: EffectiveSubscriptionView;
};

export type GetUserSubscription = Record<string, never>;

export type UpdateUserSubscriptionRequest = {
  userId?: string;
  active?: boolean;
  planKey?: string;
  featureOverrides?: Record<string, boolean>;
  couchQuotaOverrideBytes?: string;
  s3QuotaOverrideBytes?: string;
};

export type UpdateUserSubscriptionResponse = {
  effective?: EffectiveSubscriptionView;
};

export type UpdateUserSubscription = Record<string, never>;

export class AdminSubscriptionsAPI {
  static ListSubscriptionPlans(this:void, req: ListSubscriptionPlansRequest, initReq?: fm.InitReq): Promise<ListSubscriptionPlansResponse> {
    return fm.fetchRequest<ListSubscriptionPlansResponse>(`/api/admin_subscriptions/plans`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static GetUserSubscription(this:void, req: GetUserSubscriptionRequest, initReq?: fm.InitReq): Promise<GetUserSubscriptionResponse> {
    return fm.fetchRequest<GetUserSubscriptionResponse>(`/api/admin_subscriptions/get`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static UpdateUserSubscription(this:void, req: UpdateUserSubscriptionRequest, initReq?: fm.InitReq): Promise<UpdateUserSubscriptionResponse> {
    return fm.fetchRequest<UpdateUserSubscriptionResponse>(`/api/admin_subscriptions/update`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
}