package tract_templates

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/repository/pg/pg_err"
	"github.com/ruf-dev/artel/internal/repository/pg/tx_manager"
	"go.redsock.ru/rerrors"
)

// Repo is the pure-DB layer for published tract templates — immutable snapshots of a tract's
// definition, browsable/copyable by any user (see migration 053, domain.TractTemplate).
type Repo struct {
	q  *artel_q.Queries
	tx tx_manager.TxManager
}

func New(q *artel_q.Queries, tx tx_manager.TxManager) *Repo {
	return &Repo{q: q, tx: tx}
}

func (r *Repo) Create(ctx context.Context, template domain.TractTemplate) (domain.TractTemplate, error) {
	defJSON, err := json.Marshal(template.Definition)
	if err != nil {
		return domain.TractTemplate{}, rerrors.Wrap(err, "error marshaling tract template definition")
	}

	sourceTractId := uuid.NullUUID{UUID: template.SourceTractUuid, Valid: template.SourceTractUuid != uuid.Nil}

	params := artel_q.InsertTractTemplateParams{
		SourceTractID: sourceTractId,
		OwnerID:       template.OwnerUuid,
		Name:          template.Name,
		Description:   template.Description,
		Definition:    defJSON,
		Category:      template.Category,
	}

	row, err := r.q.InsertTractTemplate(ctx, params)
	if err != nil {
		return domain.TractTemplate{}, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error inserting tract template")
	}

	return tractTemplateToDomain(row)
}

func (r *Repo) Get(ctx context.Context, id uuid.UUID) (sql.Null[domain.TractTemplate], error) {
	row, err := r.q.GetTractTemplate(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.Null[domain.TractTemplate]{}, nil
		}

		return sql.Null[domain.TractTemplate]{}, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error getting tract template")
	}

	template, err := tractTemplateToDomain(row)
	if err != nil {
		return sql.Null[domain.TractTemplate]{}, err
	}

	result := sql.Null[domain.TractTemplate]{V: template, Valid: true}

	return result, nil
}

// ListAll returns published templates ordered by install_count desc, published_at desc.
// category == "" means no filter.
func (r *Repo) ListAll(ctx context.Context, category string) ([]domain.TractTemplate, error) {
	rows, err := r.q.ListTractTemplates(ctx, category)
	if err != nil {
		return nil, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error listing tract templates")
	}

	templates := make([]domain.TractTemplate, len(rows))

	for i, row := range rows {
		template, convErr := tractTemplateToDomain(row)
		if convErr != nil {
			return nil, convErr
		}

		templates[i] = template
	}

	return templates, nil
}

func (r *Repo) ListByOwner(ctx context.Context, ownerUuid uuid.UUID) ([]domain.TractTemplate, error) {
	rows, err := r.q.ListTractTemplatesByOwner(ctx, ownerUuid)
	if err != nil {
		return nil, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error listing tract templates by owner")
	}

	templates := make([]domain.TractTemplate, len(rows))

	for i, row := range rows {
		template, convErr := tractTemplateToDomain(row)
		if convErr != nil {
			return nil, convErr
		}

		templates[i] = template
	}

	return templates, nil
}

func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteTractTemplate(ctx, id)
	if err != nil {
		return rerrors.Wrap(pg_err.UnwrapPgErr(err), "error deleting tract template")
	}

	return nil
}

func (r *Repo) IncrementInstallCount(ctx context.Context, id uuid.UUID) error {
	err := r.q.IncrementTractTemplateInstallCount(ctx, id)
	if err != nil {
		return rerrors.Wrap(pg_err.UnwrapPgErr(err), "error incrementing tract template install count")
	}

	return nil
}

func tractTemplateToDomain(row artel_q.TractTemplate) (domain.TractTemplate, error) {
	var def domain.TractDefinition

	err := json.Unmarshal(row.Definition, &def)
	if err != nil {
		return domain.TractTemplate{}, rerrors.Wrap(err, "error unmarshaling tract template definition")
	}

	template := domain.TractTemplate{
		Uuid:         row.ID,
		OwnerUuid:    row.OwnerID,
		Name:         row.Name,
		Description:  row.Description,
		Definition:   def,
		Category:     row.Category,
		InstallCount: int(row.InstallCount),
		PublishedAt:  row.PublishedAt,
		UpdatedAt:    row.UpdatedAt,
	}

	if row.SourceTractID.Valid {
		template.SourceTractUuid = row.SourceTractID.UUID
	}

	return template, nil
}
