package external_connections_api

import (
	"encoding/json"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/domain"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
)

var providerToProto = map[string]pb.ExternalProvider{
	domain.ProviderGoogleSheets: pb.ExternalProvider_EXTERNAL_PROVIDER_GOOGLE_SHEETS,
	domain.ProviderTrello:       pb.ExternalProvider_EXTERNAL_PROVIDER_TRELLO,
	domain.ProviderMiro:         pb.ExternalProvider_EXTERNAL_PROVIDER_MIRO,
	domain.ProviderEmail:        pb.ExternalProvider_EXTERNAL_PROVIDER_EMAIL,
	domain.ProviderGitlab:       pb.ExternalProvider_EXTERNAL_PROVIDER_GITLAB,
}

func ConnectionToProto(m domain.ExternalConnectionMeta) *pb.ExternalConnectionInfo {
	info := &pb.ExternalConnectionInfo{
		Id:        m.Uuid.String(),
		Provider:  providerToProto[m.Provider],
		CreatedAt: m.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
		UpdatedAt: m.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}

	switch m.ProviderType {
	case artel_q.ExternalProviderTypeGoogleOauth:
		info.Details = googleDetails(m.Metadata)
	default:
		info.Details = genericDetails(m.Metadata)
	}

	return info
}

func googleDetails(raw json.RawMessage) *pb.ExternalConnectionInfo_Google {
	var meta domain.GoogleConnectionMeta
	if raw != nil {
		_ = json.Unmarshal(raw, &meta)
	}
	googleInfo := &pb.GoogleConnectionInfo{
		Email:  meta.Email,
		Scopes: meta.Scopes,
	}
	return &pb.ExternalConnectionInfo_Google{Google: googleInfo}
}

func genericDetails(raw json.RawMessage) *pb.ExternalConnectionInfo_Generic {
	fields := make(map[string]string)
	if raw != nil {
		var parsed map[string]any
		if err := json.Unmarshal(raw, &parsed); err == nil {
			for k, v := range parsed {
				if s, ok := v.(string); ok {
					fields[k] = s
				}
			}
		}
	}
	generic := &pb.GenericConnection{Fields: fields}
	return &pb.ExternalConnectionInfo_Generic{Generic: generic}
}
