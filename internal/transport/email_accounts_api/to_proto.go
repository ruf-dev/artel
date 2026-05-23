package email_accounts_api

import (
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/domain"
)

func toProto(a domain.EmailAccount) *pb.EmailAccountInfo {
	return &pb.EmailAccountInfo{
		Id:        a.Uuid.String(),
		Email:     a.Email,
		ImapHost:  a.ImapHost,
		ImapPort:  int32(a.ImapPort),
		SmtpHost:  a.SmtpHost,
		SmtpPort:  int32(a.SmtpPort),
		CreatedAt: a.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}
}
