package notes_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (n *NotesImpl) CheckImportConflicts(
	ctx context.Context, req *pb.CheckImportConflicts_Request,
) (*pb.CheckImportConflicts_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "error parsing vault id")
	}

	conflicts, err := n.noteSvc.CheckImportConflicts(ctx, vaultID, req.DestPath, req.ZipData)
	if err != nil {
		return nil, rerrors.Wrap(err, "error checking import conflicts")
	}

	resp := &pb.CheckImportConflicts_Response{ConflictingPaths: conflicts}

	return resp, nil
}
