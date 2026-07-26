// Package dockerhosts is the pure-DB layer for the admin-managed pool of Docker daemons backing
// per-vault workbench containers — see docs/workbench/02_docker_topology.md. Unlike
// couchinstances/s3instances there's no single credential blob; instead there are three optional
// TLS/mTLS fields (migrations/062_docker_hosts_tls.sql) for the remote-daemon case, so this repo
// carries an encryptor just like s3instances does.
package dockerhosts

import (
	"context"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/clients/sqldb"
	"github.com/ruf-dev/artel/internal/cryptoutil"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/repository/pg/pg_err"
	"go.redsock.ru/rerrors"
)

type Repo struct {
	q         *artel_q.Queries
	encryptor cryptoutil.Encryptor
}

func New(db sqldb.DB, encryptor cryptoutil.Encryptor) *Repo {
	return &Repo{
		q:         artel_q.New(db),
		encryptor: encryptor,
	}
}

func (r *Repo) Register(ctx context.Context, url, caCert, clientCert, clientKey string) (uuid.UUID, error) {
	caCertEnc, err := r.encryptCert(caCert)
	if err != nil {
		return uuid.UUID{}, rerrors.Wrap(err, "error encrypting ca cert")
	}

	clientCertEnc, err := r.encryptCert(clientCert)
	if err != nil {
		return uuid.UUID{}, rerrors.Wrap(err, "error encrypting client cert")
	}

	clientKeyEnc, err := r.encryptCert(clientKey)
	if err != nil {
		return uuid.UUID{}, rerrors.Wrap(err, "error encrypting client key")
	}

	params := artel_q.RegisterDockerHostParams{
		Url:           url,
		CaCertEnc:     caCertEnc,
		ClientCertEnc: clientCertEnc,
		ClientKeyEnc:  clientKeyEnc,
	}

	id, err := r.q.RegisterDockerHost(ctx, params)
	if err != nil {
		return uuid.UUID{}, pg_err.UnwrapPgErr(err)
	}

	return id, nil
}

func (r *Repo) Get(ctx context.Context, id uuid.UUID) (domain.DockerHost, error) {
	row, err := r.q.GetDockerHost(ctx, id)
	if err != nil {
		return domain.DockerHost{}, pg_err.UnwrapPgErr(err)
	}

	host := domain.DockerHost{
		Uuid:      row.ID,
		Url:       row.Url,
		CreatedAt: row.CreatedAt,
	}

	return host, nil
}

// GetWithCreds is like Get but also decrypts and populates CaCert/ClientCert/ClientKey — used
// only by workbench.Service.resolveClient to build a TLS-configured Docker client. Get/List stay
// creds-free for admin list/view, mirroring s3instances' GetS3Instance vs GetS3InstanceWithCreds
// split.
func (r *Repo) GetWithCreds(ctx context.Context, id uuid.UUID) (domain.DockerHost, error) {
	row, err := r.q.GetDockerHostWithCreds(ctx, id)
	if err != nil {
		return domain.DockerHost{}, pg_err.UnwrapPgErr(err)
	}

	caCert, err := r.decryptCert(row.CaCertEnc)
	if err != nil {
		return domain.DockerHost{}, rerrors.Wrap(err, "error decrypting ca cert")
	}

	clientCert, err := r.decryptCert(row.ClientCertEnc)
	if err != nil {
		return domain.DockerHost{}, rerrors.Wrap(err, "error decrypting client cert")
	}

	clientKey, err := r.decryptCert(row.ClientKeyEnc)
	if err != nil {
		return domain.DockerHost{}, rerrors.Wrap(err, "error decrypting client key")
	}

	host := domain.DockerHost{
		Uuid:       row.ID,
		Url:        row.Url,
		CreatedAt:  row.CreatedAt,
		CaCert:     caCert,
		ClientCert: clientCert,
		ClientKey:  clientKey,
	}

	return host, nil
}

func (r *Repo) List(ctx context.Context) ([]domain.DockerHost, error) {
	rows, err := r.q.ListDockerHosts(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "list docker hosts")
	}

	hosts := make([]domain.DockerHost, len(rows))
	for i, row := range rows {
		hosts[i] = domain.DockerHost{
			Uuid:      row.ID,
			Url:       row.Url,
			CreatedAt: row.CreatedAt,
		}
	}

	return hosts, nil
}

// Update patches url unconditionally. caCert/clientCert/clientKey are three-way patch pointers:
// nil leaves the stored cert untouched, a pointer to "" clears it (stored as NULL), a pointer to
// a non-empty PEM re-encrypts and stores it.
func (r *Repo) Update(ctx context.Context, id uuid.UUID, url string, caCert, clientCert, clientKey *string) error {
	caCertEnc, updateCaCert, err := r.encryptCertPatch(caCert)
	if err != nil {
		return rerrors.Wrap(err, "error encrypting ca cert")
	}

	clientCertEnc, updateClientCert, err := r.encryptCertPatch(clientCert)
	if err != nil {
		return rerrors.Wrap(err, "error encrypting client cert")
	}

	clientKeyEnc, updateClientKey, err := r.encryptCertPatch(clientKey)
	if err != nil {
		return rerrors.Wrap(err, "error encrypting client key")
	}

	params := artel_q.UpdateDockerHostParams{
		ID:               id,
		Url:              url,
		UpdateCaCert:     updateCaCert,
		CaCertEnc:        caCertEnc,
		UpdateClientCert: updateClientCert,
		ClientCertEnc:    clientCertEnc,
		UpdateClientKey:  updateClientKey,
		ClientKeyEnc:     clientKeyEnc,
	}

	err = r.q.UpdateDockerHost(ctx, params)
	if err != nil {
		return rerrors.Wrap(err, "update docker host")
	}

	return nil
}

func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteDockerHost(ctx, id)
	if err != nil {
		return rerrors.Wrap(err, "delete docker host")
	}

	return nil
}

func (r *Repo) Exists(ctx context.Context) (bool, error) {
	exists, err := r.q.DockerHostExists(ctx)
	if err != nil {
		return false, rerrors.Wrap(err, "check docker host existence")
	}

	return exists, nil
}

// PickLeastLoaded returns the docker_hosts row currently backing the fewest live ('created' or
// 'running') workbenches, ties broken by oldest first. Returns user_errors.NotFound (via
// pg_err.UnwrapPgErr) when no docker hosts are registered at all.
func (r *Repo) PickLeastLoaded(ctx context.Context) (domain.DockerHost, error) {
	row, err := r.q.PickLeastLoadedDockerHost(ctx)
	if err != nil {
		return domain.DockerHost{}, pg_err.UnwrapPgErr(err)
	}

	host := domain.DockerHost{
		Uuid:      row.ID,
		Url:       row.Url,
		CreatedAt: row.CreatedAt,
	}

	return host, nil
}

func (r *Repo) WithTx(tx sqldb.DB) repository.DockerHosts {
	return New(tx, r.encryptor)
}

// encryptCert encrypts plain for storage, but leaves an unset cert (empty string) as NULL rather
// than encrypting an empty string into a non-null blob — "no TLS configured" must stay
// distinguishable from "configured with an empty value" at the DB layer.
func (r *Repo) encryptCert(plain string) ([]byte, error) {
	if plain == "" {
		return nil, nil
	}

	enc, err := r.encryptor.Encrypt([]byte(plain))
	if err != nil {
		return nil, err
	}

	return enc, nil
}

// encryptCertPatch converts an Update three-way patch pointer (nil = don't touch, pointer to ""
// = clear, pointer to a PEM = set) into the (encrypted bytes, shouldTouch) pair
// UpdateDockerHost's SQL CASE WHEN branches on.
func (r *Repo) encryptCertPatch(patch *string) (enc []byte, touch bool, err error) {
	if patch == nil {
		return nil, false, nil
	}

	if *patch == "" {
		return nil, true, nil
	}

	enc, err = r.encryptor.Encrypt([]byte(*patch))
	if err != nil {
		return nil, true, err
	}

	return enc, true, nil
}

// decryptCert decrypts enc, treating an empty/NULL blob as "no cert configured" rather than
// attempting to decrypt zero bytes.
func (r *Repo) decryptCert(enc []byte) (string, error) {
	if len(enc) == 0 {
		return "", nil
	}

	plain, err := r.encryptor.Decrypt(enc)
	if err != nil {
		return "", err
	}

	return string(plain), nil
}
