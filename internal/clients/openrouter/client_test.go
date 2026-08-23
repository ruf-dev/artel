package openrouter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// newTestClient points Client at a local httptest server instead of the real
// https://openrouter.ai host, by swapping baseURL for the duration of the test.
func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client := New("test-api-key")
	client.baseURL = server.URL
	client.httpClient = server.Client()

	return client
}

func TestGetKeyInfo_UnlimitedKey(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/key", r.URL.Path)
		require.Equal(t, "Bearer test-api-key", r.Header.Get("Authorization"))

		body := `{
			"data": {
				"label": "sk-or-...abcd",
				"limit": null,
				"limit_reset": null,
				"limit_remaining": null,
				"include_byok_in_limit": false,
				"usage": 1.5,
				"usage_daily": 0.5,
				"usage_weekly": 1.0,
				"usage_monthly": 1.5,
				"byok_usage": 0,
				"byok_usage_daily": 0,
				"byok_usage_weekly": 0,
				"byok_usage_monthly": 0,
				"is_free_tier": false
			}
		}`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(body))
		require.NoError(t, err)
	}

	client := newTestClient(t, handler)

	info, err := client.GetKeyInfo(context.Background())
	require.NoError(t, err)

	require.True(t, info.LimitUnlimited)
	require.Zero(t, info.Limit)
	require.Zero(t, info.LimitRemaining)
	require.InDelta(t, 1.5, info.Usage, 0.0001)
	require.InDelta(t, 0.5, info.UsageDaily, 0.0001)
	require.InDelta(t, 1.0, info.UsageWeekly, 0.0001)
	require.InDelta(t, 1.5, info.UsageMonthly, 0.0001)
	require.False(t, info.IsFreeTier)
}

func TestGetKeyInfo_LimitedKey(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		body := `{
			"data": {
				"label": "sk-or-...wxyz",
				"limit": 100,
				"limit_reset": "monthly",
				"limit_remaining": 42.5,
				"include_byok_in_limit": true,
				"usage": 57.5,
				"usage_daily": 2,
				"usage_weekly": 10,
				"usage_monthly": 57.5,
				"byok_usage": 3,
				"byok_usage_daily": 0.1,
				"byok_usage_weekly": 1,
				"byok_usage_monthly": 3,
				"is_free_tier": true
			}
		}`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(body))
		require.NoError(t, err)
	}

	client := newTestClient(t, handler)

	info, err := client.GetKeyInfo(context.Background())
	require.NoError(t, err)

	require.False(t, info.LimitUnlimited)
	require.InDelta(t, 100, info.Limit, 0.0001)
	require.InDelta(t, 42.5, info.LimitRemaining, 0.0001)
	require.True(t, info.IncludeByokInLimit)
	require.InDelta(t, 57.5, info.Usage, 0.0001)
	require.InDelta(t, 3.0, info.ByokUsage, 0.0001)
	require.True(t, info.IsFreeTier)
}

func TestGetKeyInfo_UnauthorizedRejected(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)

		_, err := w.Write([]byte(`{"error": "no auth credentials found"}`))
		require.NoError(t, err)
	}

	client := newTestClient(t, handler)

	_, err := client.GetKeyInfo(context.Background())
	require.Error(t, err)

	statusCode, ok := StatusCode(err)
	require.True(t, ok)
	require.Equal(t, http.StatusUnauthorized, statusCode)
}
