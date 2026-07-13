package notes_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service/v1/notes"
	"go.redsock.ru/rerrors"
)

func importActionFromPb(action pb.ImportConflictAction) notes.ImportAction {
	switch action {
	case pb.ImportConflictAction_OVERWRITE:
		return notes.ImportActionOverwrite
	case pb.ImportConflictAction_RENAME:
		return notes.ImportActionRename
	default:
		return notes.ImportActionSkip
	}
}

func importResolutionsFromPb(resolutions []*pb.ImportResolution) []notes.ImportResolution {
	result := make([]notes.ImportResolution, 0, len(resolutions))

	for _, r := range resolutions {
		resolution := notes.ImportResolution{
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
