package notes_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/domain"
	"go.redsock.ru/rerrors"
)

func importActionFromPb(action pb.ImportConflictAction) domain.ImportAction {
	switch action {
	case pb.ImportConflictAction_OVERWRITE:
		return domain.ImportActionOverwrite
	case pb.ImportConflictAction_RENAME:
		return domain.ImportActionRename
	default:
		return domain.ImportActionSkip
	}
}

func importResolutionsFromPb(resolutions []*pb.ImportResolution) []domain.ImportResolution {
	result := make([]domain.ImportResolution, 0, len(resolutions))

	for _, r := range resolutions {
		resolution := domain.ImportResolution{
			Path:     r.Path,
			Action:   importActionFromPb(r.Action),
			RenameTo: r.RenameTo,
		}
		result = append(result, resolution)
	}

	return result
}

func (n *NotesImpl) CommitImport(ctx context.Context, req *pb.CommitImport_Request) (*pb.CommitImport_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "error parsing vault id")
	}

	resolutions := importResolutionsFromPb(req.Resolutions)

	imported, skipped, err := n.noteSvc.CommitImport(ctx, vaultID, req.DestPath, req.ZipData, resolutions)
	if err != nil {
		return nil, rerrors.Wrap(err, "error committing import")
	}

	resp := &pb.CommitImport_Response{
		ImportedCount: int32(imported),
		SkippedCount:  int32(skipped),
	}

	return resp, nil
}
