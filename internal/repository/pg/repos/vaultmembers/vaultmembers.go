package vaultmembers

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/clients/sqldb"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
)

type Repo struct {
	q *artel_q.Queries
}

func New(db sqldb.DB) *Repo {
	return &Repo{
		q: artel_q.New(db),
	}
}

func (r *Repo) Add(ctx context.Context, vaultUuid, userUuid uuid.UUID, role string) error {
	err := r.q.AddVaultMember(ctx, artel_q.AddVaultMemberParams{
		VaultID: vaultUuid,
		UserID:  userUuid,
		Role:    role,
	})
	if err != nil {
		return rerrors.Wrap(err, "failed to add vault member")
	}
	return nil
}

func (r *Repo) Remove(ctx context.Context, vaultUuid, userUuid uuid.UUID) error {
	err := r.q.RemoveVaultMember(ctx, artel_q.RemoveVaultMemberParams{
		VaultID: vaultUuid,
		UserID:  userUuid,
	})
	if err != nil {
		return rerrors.Wrap(err, "failed to remove vault member")
	}
	return nil
}

func (r *Repo) Get(ctx context.Context, vaultUuid, userUuid uuid.UUID) (domain.VaultMember, error) {
	member, err := r.q.GetVaultMembership(ctx, artel_q.GetVaultMembershipParams{
		VaultID: vaultUuid,
		UserID:  userUuid,
	})
	if err != nil {
		return domain.VaultMember{}, rerrors.Wrap(err, "get vault membership")
	}

	return domain.VaultMember{
		Uuid:      member.ID,
		VaultUuid: member.VaultID,
		UserUuid:  member.UserID,
		Role:      member.Role,
		CreatedAt: member.CreatedAt,
	}, nil
}

func (r *Repo) ListByVault(ctx context.Context, vaultUuid uuid.UUID) ([]domain.VaultMember, error) {
	members, err := r.q.ListVaultMembers(ctx, vaultUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "list vault members")
	}

	result := make([]domain.VaultMember, len(members))
	for i, m := range members {
		result[i] = domain.VaultMember{
			Uuid:      m.ID,
			VaultUuid: m.VaultID,
			UserUuid:  m.UserID,
			Role:      m.Role,
			CreatedAt: m.CreatedAt,
		}
	}
	return result, nil
}

func (r *Repo) WithTx(tx sqldb.DB) repository.VaultMembers {
	return New(tx)
}
