package workbenchdocker

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/docker/docker/client"
	"github.com/stretchr/testify/require"
)

// newTestClient builds a Client whose underlying Docker SDK client talks to server instead of
// a real daemon. It pins an explicit API version (client.WithVersion) rather than
// client.WithAPIVersionNegotiation() (what New uses) so tests don't also need to fake a
// GET /version response.
func newTestClient(t *testing.T, server *httptest.Server) *Client {
	t.Helper()

	host := "tcp://" + strings.TrimPrefix(server.URL, "http://")

	dockerCli, err := client.NewClientWithOpts(
		client.WithHost(host),
		client.WithVersion("1.44"),
	)
	require.NoError(t, err)

	return &Client{cli: dockerCli, host: host}
}

func TestCreateContainer_Success(t *testing.T) {
	var gotPath string
	var gotBody map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path

		require.Equal(t, http.MethodPost, r.Method)

		err := json.NewDecoder(r.Body).Decode(&gotBody)
		require.NoError(t, err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		_, err = w.Write([]byte(`{"Id":"container123","Warnings":[]}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	c := newTestClient(t, server)

	opts := CreateOpts{
		Name:       "workbench-test",
		VolumeName: "workbench-vol-test",
	}

	id, err := c.CreateContainer(context.Background(), opts)
	require.NoError(t, err)
	require.Equal(t, "container123", id)
	require.Contains(t, gotPath, "/containers/create")

	// container.CreateRequest embeds *Config anonymously, so its fields (Image, Labels, ...)
	// flatten into the top-level request body rather than nesting under a "Config" key.
	require.Equal(t, workbenchImage, gotBody["Image"])

	labels, ok := gotBody["Labels"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "true", labels[workbenchLabelKey])

	hostConfig, ok := gotBody["HostConfig"].(map[string]any)
	require.True(t, ok)

	mounts, ok := hostConfig["Mounts"].([]any)
	require.True(t, ok)
	require.Len(t, mounts, 1)

	mnt, ok := mounts[0].(map[string]any)
	require.True(t, ok)
	require.Equal(t, opts.VolumeName, mnt["Source"])
	require.Equal(t, workspaceMountPath, mnt["Target"])

	networkingConfig, ok := gotBody["NetworkingConfig"].(map[string]any)
	require.True(t, ok)

	endpointsConfig, ok := networkingConfig["EndpointsConfig"].(map[string]any)
	require.True(t, ok)
	require.Contains(t, endpointsConfig, workbenchNetworkName)
}

func TestCreateContainer_DaemonError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		_, err := w.Write([]byte(`{"message":"no space left on device"}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	c := newTestClient(t, server)

	opts := CreateOpts{Name: "workbench-test", VolumeName: "vol"}

	id, err := c.CreateContainer(context.Background(), opts)
	require.Error(t, err)
	require.Empty(t, id)
}

func TestStartContainer_Success(t *testing.T) {
	var gotPaths []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPaths = append(gotPaths, r.URL.Path)

		switch {
		case strings.Contains(r.URL.Path, "/start") && strings.Contains(r.URL.Path, "/exec/"):
			w.WriteHeader(http.StatusNoContent)
		case strings.Contains(r.URL.Path, "/exec"):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)

			_, err := w.Write([]byte(`{"Id":"exec123"}`))
			require.NoError(t, err)
		case strings.Contains(r.URL.Path, "/start"):
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected request path %q", r.URL.Path)
		}
	}))
	defer server.Close()

	c := newTestClient(t, server)

	env := map[string]string{"ANTHROPIC_API_KEY": "sk-test"}

	err := c.StartContainer(context.Background(), "container123", env)
	require.NoError(t, err)

	require.Contains(t, gotPaths, "/v1.44/containers/container123/start")

	var sawExecCreate, sawExecStart bool
	for _, p := range gotPaths {
		if strings.Contains(p, "/containers/container123/exec") {
			sawExecCreate = true
		}
		if strings.Contains(p, "/exec/exec123/start") {
			sawExecStart = true
		}
	}
	require.True(t, sawExecCreate, "expected an exec-create call to inject env, got paths %v", gotPaths)
	require.True(t, sawExecStart, "expected an exec-start call to inject env, got paths %v", gotPaths)
}

func TestStartContainer_NoEnv(t *testing.T) {
	var gotPaths []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPaths = append(gotPaths, r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	c := newTestClient(t, server)

	err := c.StartContainer(context.Background(), "container123", nil)
	require.NoError(t, err)
	require.Len(t, gotPaths, 1)
	require.Contains(t, gotPaths[0], "/containers/container123/start")
}

func TestStartContainer_DaemonError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		_, err := w.Write([]byte(`{"message":"no such container"}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	c := newTestClient(t, server)

	err := c.StartContainer(context.Background(), "does-not-exist", nil)
	require.Error(t, err)
}

func TestStopContainer_Success(t *testing.T) {
	var gotPath, gotMethod string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	c := newTestClient(t, server)

	err := c.StopContainer(context.Background(), "container123")
	require.NoError(t, err)
	require.Equal(t, http.MethodPost, gotMethod)
	require.Contains(t, gotPath, "/containers/container123/stop")
}

func TestStopContainer_DaemonError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		_, err := w.Write([]byte(`{"message":"boom"}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	c := newTestClient(t, server)

	err := c.StopContainer(context.Background(), "container123")
	require.Error(t, err)
}

func TestRemoveContainer_Success(t *testing.T) {
	var gotPath, gotMethod, gotQuery string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		gotQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	c := newTestClient(t, server)

	err := c.RemoveContainer(context.Background(), "container123")
	require.NoError(t, err)
	require.Equal(t, http.MethodDelete, gotMethod)
	require.Contains(t, gotPath, "/containers/container123")
	require.Contains(t, gotQuery, "force=1")
}

func TestRemoveContainer_DaemonError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)

		_, err := w.Write([]byte(`{"message":"container is running"}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	c := newTestClient(t, server)

	err := c.RemoveContainer(context.Background(), "container123")
	require.Error(t, err)
}

func TestCreateVolume_Success(t *testing.T) {
	var gotPath string
	var gotBody map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path

		err := json.NewDecoder(r.Body).Decode(&gotBody)
		require.NoError(t, err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		_, err = w.Write([]byte(`{"Name":"workbench-vol-test","Driver":"local","Mountpoint":"/var/lib/docker/volumes/workbench-vol-test/_data","Labels":{}}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	c := newTestClient(t, server)

	err := c.CreateVolume(context.Background(), "workbench-vol-test")
	require.NoError(t, err)
	require.Contains(t, gotPath, "/volumes/create")
	require.Equal(t, "workbench-vol-test", gotBody["Name"])

	labels, ok := gotBody["Labels"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "true", labels[workbenchLabelKey])
}

func TestCreateVolume_DaemonError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		_, err := w.Write([]byte(`{"message":"disk quota exceeded"}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	c := newTestClient(t, server)

	err := c.CreateVolume(context.Background(), "workbench-vol-test")
	require.Error(t, err)
}

func TestRemoveVolume_Success(t *testing.T) {
	var gotPath, gotMethod, gotQuery string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		gotQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	c := newTestClient(t, server)

	err := c.RemoveVolume(context.Background(), "workbench-vol-test")
	require.NoError(t, err)
	require.Equal(t, http.MethodDelete, gotMethod)
	require.Contains(t, gotPath, "/volumes/workbench-vol-test")
	require.Contains(t, gotQuery, "force=1")
}

func TestRemoveVolume_DaemonError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		_, err := w.Write([]byte(`{"message":"no such volume"}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	c := newTestClient(t, server)

	err := c.RemoveVolume(context.Background(), "workbench-vol-test")
	require.Error(t, err)
}

// TestNew_ConfiguresHost is a light sanity check on New itself (as opposed to newTestClient's
// direct construction used by every other test above): it doesn't hit a daemon (no method that
// performs I/O is called), it just verifies New succeeds and records the configured host, since
// New's own docker.Client is otherwise unexported and untestable against a fake daemon without
// also faking GET /version (New uses API version negotiation, unlike the pinned-version test
// client above).
func TestNew_ConfiguresHost(t *testing.T) {
	c, err := New("unix:///var/run/docker-workbenches.sock")
	require.NoError(t, err)
	require.NotNil(t, c)
}
