package couchdb

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
)

// CredentialSource loads per-vault CouchDB credentials from persistent storage.
// Satisfied by repository.CouchCredentials.
type CredentialSource interface {
	Load(ctx context.Context, vaultID uuid.UUID) (domain.CouchCred, error)
}

// Pool caches CouchDB clients keyed by BaseURL. Credentials are fetched from
// Postgres on the first call for a given vault and reused on subsequent calls.
// The defaultClient is the admin-level client used for provisioning new vaults.
type Pool struct {
	mu      sync.RWMutex
	clients map[string]*Client
	creds   CredentialSource
}

func NewPool(creds CredentialSource) *Pool {
	return &Pool{
		clients: make(map[string]*Client),
		creds:   creds,
	}
}

func (p *Pool) Get(ctx context.Context, vaultID uuid.UUID) (*Client, error) {
	cred, err := p.creds.Load(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "load credentials")
	}

	p.mu.RLock()
	client, ok := p.clients[cred.Host]
	p.mu.RUnlock()
	if ok {
		return client, nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if client, ok = p.clients[cred.Host]; ok {
		return client, nil
	}

	client = New(Config{
		BaseURL:  cred.Host,
		User:     cred.Username,
		Password: string(cred.Password),
	})
	p.clients[cred.Host] = client

	return client, nil
}
