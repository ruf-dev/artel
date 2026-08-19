// Package workbenchdocker is a thin typed wrapper around the Docker Engine SDK
// (github.com/docker/docker/client), used to create/start/stop/remove the containers and
// volumes that back a single workbench (a per-user sandbox running two independent ways to drive
// the `claude` CLI: the in-container chat bridge, deploy/workbench/bridge, which drives it one
// turn at a time, and a tmux session (see tmux_tabs.go) exposing a real interactive `claude` TUI
// per tab through ttyd).
//
// This package talks to whichever daemon URL it's constructed with — one client per call,
// pointed at whichever docker_hosts row a given workbench is assigned to (resolved by
// internal/service/v1/workbench/workbench.go's resolveClient), not a single startup-time config
// value: a docker host is expected to be a dedicated second dockerd process, not the daemon
// Artel's own containers run on. It intentionally does not create the `workbench-net` network
// itself (assumed to pre-exist on the configured daemon) and does not expose any inbound port on
// the containers it creates.
package workbenchdocker

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/url"

	"github.com/docker/docker/client"

	"go.redsock.ru/rerrors"
)

const (
	// workbenchNetworkName is the dedicated, isolated Docker network workbench containers are
	// attached to. Not created by this client — assumed to pre-exist on the configured daemon.
	workbenchNetworkName = "workbench-net"

	// workspaceMountPath is the fixed in-container path the workbench's named volume is
	// mounted at.
	workspaceMountPath = "/workspace"

	// workbenchLabelKey/workbenchLabelValue tag every container/volume this client creates,
	// for operational visibility even on a dedicated daemon.
	workbenchLabelKey   = "artel.workbench"
	workbenchLabelValue = "true"

	// Resource limits, hardcoded at create time for the prototype.
	workbenchCpuLimitNanoCpus = 1_000_000_000          // 1 CPU
	workbenchMemLimitBytes    = 2 * 1024 * 1024 * 1024 // 2GB
)

// Client wraps the Docker SDK client with the narrow surface a workbench's container/volume
// lifecycle needs.
type Client struct {
	cli  *client.Client
	host string
}

// TLSConfig carries the decrypted, in-memory PEM-encoded TLS client credentials for connecting
// to a remote Docker daemon over TLS/mTLS — the docker_hosts.ca_cert_enc/client_cert_enc/
// client_key_enc columns (migrations/062_docker_hosts_tls.sql), decrypted by
// dockerhosts.Repo.GetWithCreds. All three fields empty means "no TLS" — see New.
type TLSConfig struct {
	CaCert     string
	ClientCert string
	ClientKey  string
}

// enabled reports whether any TLS material was provided. A partially-filled TLSConfig (e.g. a
// CA cert with no client cert/key) is treated as enabled too — X509KeyPair will fail loudly on
// the missing half rather than silently falling back to no TLS.
func (t TLSConfig) enabled() bool {
	return t.CaCert != "" || t.ClientCert != "" || t.ClientKey != ""
}

// New constructs a Client talking to the Docker daemon at host (e.g.
// "unix:///var/run/docker-workbenches.sock", "tcp://host:2376" for a local/insecure daemon, or
// "tcp://host:2376" with a populated tlsCfg for a remote mTLS-secured daemon). API version
// negotiation is enabled so the client stays compatible with the daemon regardless of exactly
// which API version it speaks.
//
// An empty tlsCfg reproduces exactly the pre-TLS behavior (plain client.WithHost, no custom HTTP
// client) — this must not regress the local unix-socket/dind path.
func New(host string, tlsCfg TLSConfig) (*Client, error) {
	opts := []client.Opt{}

	if tlsCfg.enabled() {
		httpClient, err := newTLSHTTPClient(tlsCfg)
		if err != nil {
			return nil, rerrors.Wrap(err, "error building TLS http client")
		}

		// WithHTTPClient must be set before WithHost — WithHost type-asserts the client's
		// transport to configure TLS-adjacent options, so the custom *http.Client needs to
		// already be in place when it runs.
		opts = append(opts, client.WithHTTPClient(httpClient))
	}

	opts = append(opts, client.WithHost(host), client.WithAPIVersionNegotiation())

	cli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		return nil, rerrors.Wrap(err, "creating docker client")
	}

	c := &Client{
		cli:  cli,
		host: host,
	}

	return c, nil
}

// daemonHost returns the hostname portion of c.host — the address other Docker API calls on this
// Client already reach the daemon at, and (per ContainerAddress's doc comment) the address a
// published container port is reachable on too. A "unix://" (or Windows "npipe://") socket means
// the daemon is a local process on this same host, so the loopback address is used instead of a
// meaningless socket path.
func (c *Client) daemonHost() (string, error) {
	parsed, err := url.Parse(c.host)
	if err != nil {
		return "", rerrors.Wrap(err, "parsing docker host")
	}

	switch parsed.Scheme {
	case "unix", "npipe":
		return "127.0.0.1", nil
	}

	hostname := parsed.Hostname()
	if hostname == "" {
		return "", rerrors.New("docker host has no hostname: " + c.host)
	}

	return hostname, nil
}

// newTLSHTTPClient builds an *http.Client whose transport presents cfg's client certificate and
// trusts cfg's CA — the Docker SDK's client.WithTLSClientConfig only accepts file paths, but this
// repo passes decrypted secrets as in-memory values everywhere else, so the TLS config is built
// by hand here instead of round-tripping cert material through temp files.
func newTLSHTTPClient(cfg TLSConfig) (*http.Client, error) {
	cert, err := tls.X509KeyPair([]byte(cfg.ClientCert), []byte(cfg.ClientKey))
	if err != nil {
		return nil, rerrors.Wrap(err, "error parsing client cert/key pair")
	}

	caPool := x509.NewCertPool()

	ok := caPool.AppendCertsFromPEM([]byte(cfg.CaCert))
	if !ok {
		return nil, rerrors.New("error parsing CA certificate PEM")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caPool,
		MinVersion:   tls.VersionTLS12,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	httpClient := &http.Client{
		Transport: transport,
	}

	return httpClient, nil
}
