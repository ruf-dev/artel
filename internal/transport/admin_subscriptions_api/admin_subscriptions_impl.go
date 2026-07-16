package admin_subscriptions_api

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service"
	"go.redsock.ru/rerrors"
	"google.golang.org/grpc"
)

type AdminSubscriptionsImpl struct {
	artel_api.UnimplementedAdminSubscriptionsAPIServer
	svc service.AdminSubscriptionsService
}

func New(svc service.AdminSubscriptionsService) *AdminSubscriptionsImpl {
	return &AdminSubscriptionsImpl{svc: svc}
}

func (a *AdminSubscriptionsImpl) Register(srv grpc.ServiceRegistrar) {
	artel_api.RegisterAdminSubscriptionsAPIServer(srv, a)
}

func (a *AdminSubscriptionsImpl) Gateway(ctx context.Context, endpoint string, opts ...grpc.DialOption) (string, http.Handler) {
	gwMux := runtime.NewServeMux()

	err := artel_api.RegisterAdminSubscriptionsAPIHandlerFromEndpoint(ctx, gwMux, endpoint, opts)
	if err != nil {
		log.Error().Err(err).Msg("error registering admin subscriptions grpc-gateway handler")
	}

	return "/api/admin_subscriptions/", gwMux
}

func (a *AdminSubscriptionsImpl) ListSubscriptionPlans(
	ctx context.Context,
	_ *artel_api.ListSubscriptionPlans_Request,
) (*artel_api.ListSubscriptionPlans_Response, error) {
	plans, err := a.svc.ListPlans(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing subscription plans")
	}

	entries := make([]*artel_api.SubscriptionPlanEntry, len(plans))

	for i, plan := range plans {
		entries[i] = subscriptionPlanEntry(plan)
	}

	resp := &artel_api.ListSubscriptionPlans_Response{Plans: entries}

	return resp, nil
}

func (a *AdminSubscriptionsImpl) GetUserSubscription(
	ctx context.Context,
	req *artel_api.GetUserSubscription_Request,
) (*artel_api.GetUserSubscription_Response, error) {
	id, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, rerrors.Wrap(err, "error parsing user id")
	}

	sub, effective, err := a.svc.GetUserSubscription(ctx, id)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting user subscription")
	}

	resp := &artel_api.GetUserSubscription_Response{
		Active:                  sub.Active,
		PlanKey:                 sub.PlanKey,
		FeatureOverrides:        featureOverridesToProto(sub.FeatureOverrides),
		CouchQuotaOverrideBytes: sub.CouchQuotaOverrideBytes,
		S3QuotaOverrideBytes:    sub.S3QuotaOverrideBytes,
		Effective:               effectiveSubscriptionView(effective),
	}

	return resp, nil
}

func (a *AdminSubscriptionsImpl) UpdateUserSubscription(
	ctx context.Context,
	req *artel_api.UpdateUserSubscription_Request,
) (*artel_api.UpdateUserSubscription_Response, error) {
	id, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, rerrors.Wrap(err, "error parsing user id")
	}

	sub := domain.Subscription{
		UserUuid:                id,
		Active:                  req.Active,
		PlanKey:                 req.PlanKey,
		FeatureOverrides:        featureOverridesFromProto(req.FeatureOverrides),
		CouchQuotaOverrideBytes: req.CouchQuotaOverrideBytes,
		S3QuotaOverrideBytes:    req.S3QuotaOverrideBytes,
	}

	effective, err := a.svc.UpdateUserSubscription(ctx, sub)
	if err != nil {
		return nil, rerrors.Wrap(err, "error updating user subscription")
	}

	resp := &artel_api.UpdateUserSubscription_Response{Effective: effectiveSubscriptionView(effective)}

	return resp, nil
}

func subscriptionPlanEntry(plan domain.SubscriptionPlan) *artel_api.SubscriptionPlanEntry {
	entry := &artel_api.SubscriptionPlanEntry{
		PlanKey:         plan.PlanKey,
		CouchQuotaBytes: plan.CouchQuotaBytes,
		S3QuotaBytes:    plan.S3QuotaBytes,
		Features:        featureSetToProto(plan.Features),
	}

	return entry
}

func effectiveSubscriptionView(effective domain.EffectiveSubscription) *artel_api.EffectiveSubscriptionView {
	view := &artel_api.EffectiveSubscriptionView{
		Active:          effective.Active,
		PlanKey:         effective.PlanKey,
		Features:        featureSetToProto(effective.Features),
		CouchQuotaBytes: effective.CouchQuotaBytes,
		S3QuotaBytes:    effective.S3QuotaBytes,
	}

	return view
}

// allFeatures enumerates every domain.SubscriptionFeature so a FeatureSet can be flattened into
// a full map for display — unlike feature_overrides, this map is never sparse.
var allFeatures = []domain.SubscriptionFeature{
	domain.FeatureEmails,
	domain.FeatureTaskTrackers,
	domain.FeatureNotes,
	domain.FeatureSpreadsheets,
	domain.FeatureConnectors,
	domain.FeatureTract,
}

func featureSetToProto(features domain.FeatureSet) map[string]bool {
	result := make(map[string]bool, len(allFeatures))

	for _, feature := range allFeatures {
		result[string(feature)] = features.Has(feature)
	}

	return result
}

func featureOverridesToProto(overrides map[domain.SubscriptionFeature]bool) map[string]bool {
	result := make(map[string]bool, len(overrides))

	for feature, value := range overrides {
		result[string(feature)] = value
	}

	return result
}

func featureOverridesFromProto(overrides map[string]bool) map[domain.SubscriptionFeature]bool {
	result := make(map[domain.SubscriptionFeature]bool, len(overrides))

	for feature, value := range overrides {
		result[domain.SubscriptionFeature(feature)] = value
	}

	return result
}
