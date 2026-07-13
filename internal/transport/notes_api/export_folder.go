package notes_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (n *NotesImpl) ExportFolder(ctx context.Context, req *pb.ExportFolder_Request) (*pb.ExportFolder_Response, error) {
	vaultID, err := uuid.Parse(req.VaultId)
	if err != nil {
		return nil, rerrors.Wrap(err, "error parsing vault id")
	}

	zipData, err := n.noteSvc.ExportFolder(ctx, vaultID, req.Path)
	if err != nil {
		return nil, rerrors.Wrap(err, "error exporting folder")
	}

	resp := &pb.ExportFolder_Response{ZipData: zipData}

	return resp, nil
}
