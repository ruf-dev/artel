package tract

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/service/v1/subscription"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTriggerTestService(
	triggers *fakeTriggersRepo,
	triggerPresets *fakeTriggerPresetsRepo,
	externalConns *fakeExternalConnsRepo,
) *Service {
	svc := New(nil, triggers, triggerPresets, externalConns, nil, nil, subscription.NewFree())

	return svc
}

func gitlabPushPreset() domain.TriggerPreset {
	return domain.TriggerPreset{
		Key:      SourceGitlabPush,
		Category: "gitlab",
		Label:    "Branch push",
		Provider: domain.ProviderGitlab,
		DefaultMatchers: domain.TriggerMatchers{
			CheckHeaders: []domain.HeaderMatcher{{Header: "X-Gitlab-Event", Equals: "Push Hook"}},
		},
	}
}

func TestCreateTrigger_ProviderLinked_NoConnection_Errors(t *testing.T) {
	triggers := newFakeTriggersRepo()
	presets := newFakeTriggerPresetsRepo()
	presets.addPreset(gitlabPushPreset())
	externalConns := newFakeExternalConnsRepo()

	svc := newTriggerTestService(triggers, presets, externalConns)

	uc := user_context.UserContext{UserUuid: uuid.New()}
	ctx := user_context.WithUserContext(context.Background(), uc)

	_, rawToken, err := svc.CreateTrigger(ctx, "gitlab trigger", "webhook", SourceGitlabPush, nil, domain.ToolSchema{})

	require.Error(t, err)
	assert.Empty(t, rawToken)
}

func TestCreateTrigger_ProviderLinked_WithConnection_LinksNoSecret(t *testing.T) {
	triggers := newFakeTriggersRepo()
	presets := newFakeTriggerPresetsRepo()
	presets.addPreset(gitlabPushPreset())
	externalConns := newFakeExternalConnsRepo()

	userUuid := uuid.New()
	conn := domain.ExternalConnection{
		Uuid:     uuid.New(),
		UserUuid: userUuid,
		Provider: domain.ProviderGitlab,
	}

	_, err := externalConns.Upsert(context.Background(), conn)
	require.NoError(t, err)

	svc := newTriggerTestService(triggers, presets, externalConns)

	uc := user_context.UserContext{UserUuid: userUuid}
	ctx := user_context.WithUserContext(context.Background(), uc)

	created, rawToken, err := svc.CreateTrigger(ctx, "gitlab trigger", "webhook", SourceGitlabPush, nil, domain.ToolSchema{})

	require.NoError(t, err)
	assert.Empty(t, rawToken)
	assert.Nil(t, created.SecretHash)
	assert.Equal(t, gitlabPushPreset().DefaultMatchers, created.Matchers)

	linkedConn, ok := triggers.providerLinks[created.Uuid]
	require.True(t, ok)
	assert.Equal(t, conn.Uuid, linkedConn)
}

func TestCreateTrigger_Generic_MintsSecretNoProviderLink(t *testing.T) {
	triggers := newFakeTriggersRepo()
	presets := newFakeTriggerPresetsRepo()
	presets.addPreset(domain.TriggerPreset{Key: SourceGeneric, Category: "generic", Label: "Generic webhook"})
	externalConns := newFakeExternalConnsRepo()

	svc := newTriggerTestService(triggers, presets, externalConns)

	uc := user_context.UserContext{UserUuid: uuid.New()}
	ctx := user_context.WithUserContext(context.Background(), uc)

	created, rawToken, err := svc.CreateTrigger(ctx, "generic trigger", "webhook", SourceGeneric, nil, domain.ToolSchema{})

	require.NoError(t, err)
	assert.NotEmpty(t, rawToken)
	assert.NotEmpty(t, created.SecretHash)

	_, linked := triggers.providerLinks[created.Uuid]
	assert.False(t, linked)
}
