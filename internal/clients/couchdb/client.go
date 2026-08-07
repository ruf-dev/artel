package couchdb

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-kivik/kivik/v4"
	kivikcouch "github.com/go-kivik/kivik/v4/couchdb"

	"github.com/ruf-dev/artel/internal/service/user_errors"
	"github.com/ruf-dev/artel/internal/utils"

	"go.redsock.ru/rerrors"
)

// enableSingleNodeMaxAttempts and enableSingleNodeRetryDelay bound the retry in enableSingleNode
// — see that function's comment for why a retry is needed at all here. This is on top of, not a
// replacement for, loggingTransport's own retry of a transient 401 (log_transport.go) — this one
// also covers enable_single_node's own 400 ("already enabled") short-circuit, which isn't a retry
// case at the transport layer.
const (
	enableSingleNodeMaxAttempts = 5
	enableSingleNodeRetryDelay  = 300 * time.Millisecond
)

type Config struct {
	BaseURL  string
	User     string
	Password string
}

type Client struct {
	kivik      *kivik.Client
	httpClient *http.Client
	baseURL    string
	user       string
	password   string
}

func New(cfg Config) (*Client, error) {
	httpClient := &http.Client{Transport: newLoggingTransport()}

	kc, err := kivik.New("couch", cfg.BaseURL,
		kivikcouch.BasicAuth(cfg.User, cfg.Password),
		kivikcouch.OptionHTTPClient(httpClient),
	)
	if err != nil {
		return nil, rerrors.Wrap(err, "creating kivik client")
	}

	c := &Client{
		kivik:      kc,
		httpClient: httpClient,
		baseURL:    cfg.BaseURL,
		user:       cfg.User,
		password:   cfg.Password,
	}

	return c, nil
}

func (c *Client) Config() Config {
	return Config{
		BaseURL:  c.baseURL,
		User:     c.user,
		Password: c.password,
	}
}

func (c *Client) Setup(ctx context.Context) error {
	err := c.enableSingleNode(ctx)
	if err != nil {
		return rerrors.Wrap(err, "enabling single node")
	}

	for _, db := range []string{"_users", "_replicator"} {
		err = c.CreateDatabase(ctx, db)
		if err != nil && !errors.Is(err, user_errors.CouchDbDatabaseAlreadyExists) {
			return rerrors.Wrap(err, "create system db "+db)
		}
	}

	return nil
}

// enableSingleNode retries on failure: when multiple callers hit /_cluster_setup concurrently
// against the same fresh CouchDB node (e.g. several test suites racing to set up the one shared
// instance at once), CouchDB can transiently reject a request with "Unauthorized" while it's still
// propagating the admin credentials this same call is in the middle of writing — a caller that
// eventually retries against the now-stabilized node succeeds. A 400 Bad Request (single node
// already enabled) is treated as success, not retried.
func (c *Client) enableSingleNode(ctx context.Context) error {
	payload := map[string]interface{}{
		"action":       "enable_single_node",
		"bind_address": "0.0.0.0",
		"username":     c.user,
		"password":     c.password,
	}

	var err error

	for attempt := 0; attempt < enableSingleNodeMaxAttempts; attempt++ {
		err = c.kivik.ClusterSetup(ctx, payload)
		if err == nil || kivik.HTTPStatus(err) == http.StatusBadRequest {
			return nil
		}

		if attempt < enableSingleNodeMaxAttempts-1 {
			time.Sleep(enableSingleNodeRetryDelay)
		}
	}

	return rerrors.Wrap(err, "cluster setup")
}

func (c *Client) GetSetupStatus(ctx context.Context) (SetupStatus, error) {
	clusterEnabled, err := c.isClusterModeEnabled(ctx)
	if err != nil {
		return SetupStatus{}, rerrors.Wrap(err, "checking cluster mode")
	}

	usersExists, err := c.kivik.DBExists(ctx, "_users")
	if err != nil {
		return SetupStatus{}, rerrors.Wrap(err, "checking _users db")
	}

	replicatorExists, err := c.kivik.DBExists(ctx, "_replicator")
	if err != nil {
		return SetupStatus{}, rerrors.Wrap(err, "checking _replicator db")
	}

	status := SetupStatus{
		ClusterModeEnabled:      clusterEnabled,
		UsersDbInitialized:      usersExists,
		ReplicatorDbInitialized: replicatorExists,
	}

	return status, nil
}

func (c *Client) isClusterModeEnabled(ctx context.Context) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/_cluster_setup", nil)
	if err != nil {
		return false, rerrors.Wrap(err, "error creating cluster status request")
	}

	req.SetBasicAuth(c.user, c.password)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, rerrors.Wrap(err, "error sending cluster status request")
	}
	defer utils.CloseWithLog(resp.Body, "error closing cluster status request")

	var body struct {
		State string `json:"state"`
	}

	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return false, rerrors.Wrap(err, "error decoding cluster status response")
	}

	return body.State != "setup_required", nil
}
