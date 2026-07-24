package tract

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestServiceWithTemplates builds a Service backed by a fakeTractTemplatesRepo, for tests
// that only exercise the template-related methods (PublishTemplate/UnpublishTemplate/
// ListTemplates/GetTemplate).
func newTestServiceWithTemplates(templates *fakeTractTemplatesRepo) *Service {
	svc := New(nil, templates, nil, nil, nil, nil, nil, nil, nil)

	return svc
}

func TestUnpublishTemplate_BuiltinNeverOwned(t *testing.T) {
	templates := newFakeTractTemplatesRepo()

	builtin := domain.TractTemplate{
		Uuid:      uuid.New(),
		OwnerUuid: uuid.Nil,
		Name:      "Create MR on featurep push",
		Category:  "Mr on push",
	}
	templates.templates[builtin.Uuid] = builtin

	svc := newTestServiceWithTemplates(templates)

	callerUuid := uuid.New()
	uc := user_context.UserContext{UserUuid: callerUuid}
	ctx := user_context.WithUserContext(context.Background(), uc)

	err := svc.UnpublishTemplate(ctx, builtin.Uuid)
	require.Error(t, err)
	assert.True(t, errors.Is(err, user_errors.TractTemplateNotOwned))

	// The template must still be present — the delete was never reached.
	_, ok := templates.templates[builtin.Uuid]
	assert.True(t, ok)
}

func TestListTemplates_IncludesBuiltins(t *testing.T) {
	templates := newFakeTractTemplatesRepo()

	builtin := domain.TractTemplate{
		Uuid:      uuid.New(),
		OwnerUuid: uuid.Nil,
		Name:      "Create MR on featurep push",
		Category:  "Mr on push",
	}
	templates.templates[builtin.Uuid] = builtin

	svc := newTestServiceWithTemplates(templates)

	callerUuid := uuid.New()
	uc := user_context.UserContext{UserUuid: callerUuid}
	ctx := user_context.WithUserContext(context.Background(), uc)

	got, err := svc.ListTemplates(ctx, "", false)
	require.NoError(t, err)

	var found bool

	for _, template := range got {
		if template.Uuid == builtin.Uuid {
			found = true
		}
	}

	assert.True(t, found, "expected ListTemplates(mineOnly=false) to include the built-in template")
}
