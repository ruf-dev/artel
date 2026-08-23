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
	domain.ProviderAnthropic:    pb.ExternalProvider_EXTERNAL_PROVIDER_ANTHROPIC,
	domain.ProviderOpenAI:       pb.ExternalProvider_EXTERNAL_PROVIDER_OPENAI,
	domain.ProviderOpenRouter:   pb.ExternalProvider_EXTERNAL_PROVIDER_OPENROUTER,
	domain.ProviderS3:           pb.ExternalProvider_EXTERNAL_PROVIDER_S3,
	domain.ProviderCouchDB:      pb.ExternalProvider_EXTERNAL_PROVIDER_COUCHDB,
	domain.ProviderPostgres:     pb.ExternalProvider_EXTERNAL_PROVIDER_POSTGRES,
}

// OpenAICompatibleProviderFromProto resolves the domain provider constant AddOpenAIConnection /
// CheckOpenAIConnection should operate on for the given request's provider field. Both RPCs are
// shared across every provider that speaks the OpenAI Chat Completions protocol, distinguished
// only by this field; EXTERNAL_PROVIDER_UNSPECIFIED preserves pre-OpenRouter callers' behavior by
// defaulting to OpenAI.
func OpenAICompatibleProviderFromProto(p pb.ExternalProvider) string {
	if p == pb.ExternalProvider_EXTERNAL_PROVIDER_OPENROUTER {
		return domain.ProviderOpenRouter
	}

	return domain.ProviderOpenAI
}

func ConnectionToProto(m domain.ExternalConnectionMeta) *pb.ExternalConnectionInfo {
	info := &pb.ExternalConnectionInfo{
		Id:           m.Uuid.String(),
		Provider:     providerToProto[m.Provider],
		ProviderName: m.Provider,
		CreatedAt:    m.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    m.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z"),
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

		err := json.Unmarshal(raw, &parsed)
		if err == nil {
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
