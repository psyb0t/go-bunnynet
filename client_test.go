package bunnynet_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/psyb0t/go-bunnynet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testAPIKey = "EXAMPLE-DO-NOT-USE"

func newTestClient(
	t *testing.T,
	handler http.HandlerFunc,
) *bunnynet.ClientWithResponses {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	auth := bunnynet.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set("AccessKey", testAPIKey)

			return nil
		},
	)

	client, err := bunnynet.NewClientWithResponses(srv.URL, auth)
	require.NoError(t, err)

	return client
}

func TestListPullZones(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/pullzone", r.URL.Path)
		assert.Equal(t, testAPIKey, r.Header.Get("AccessKey"))

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`[
			{"Id":42,"Name":"assets","OriginUrl":null,
			"MonthlyCharges":1234.567890123}
		]`))
		assert.NoError(t, err)
	})

	res, err := client.PullZonePublicIndexAllWithResponse(
		context.Background(),
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.NotNil(t, res.JSON200)
	require.Len(t, *res.JSON200, 1)

	zone := (*res.JSON200)[0]
	assert.Equal(t, int64(42), *zone.Id)
	assert.Equal(t, "assets", *zone.Name)
	assert.Nil(t, zone.OriginUrl)
	assert.InDelta(t, 1234.567890123, *zone.MonthlyCharges, 0)
}

func TestGetPullZone(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/pullzone/42", r.URL.Path)
		assert.Equal(t, "true", r.URL.Query().Get("includeCertificate"))
		assert.Equal(t, testAPIKey, r.Header.Get("AccessKey"))

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"Id":42,"Name":"assets"}`))
		assert.NoError(t, err)
	})

	res, err := client.PullZonePublicIndexWithResponse(
		context.Background(),
		42,
		&bunnynet.PullZonePublicIndexParams{IncludeCertificate: new(true)},
	)
	require.NoError(t, err)
	require.NotNil(t, res.JSON200)
	assert.Equal(t, int64(42), *res.JSON200.Id)
	assert.Equal(t, "assets", *res.JSON200.Name)
}

func TestAddPullZone(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/pullzone", r.URL.Path)
		assert.Equal(t, testAPIKey, r.Header.Get("AccessKey"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var body bunnynet.PullZoneAddModel

		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "assets", body.Name)
		w.WriteHeader(http.StatusCreated)
	})

	res, err := client.PullZonePublicAddWithResponse(
		context.Background(),
		bunnynet.PullZoneAddModel{Name: "assets"},
	)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode())
	assert.Empty(t, res.Body)
}

func TestAPIErrorStatus(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		status int
	}{
		{name: "unauthorized", status: http.StatusUnauthorized},
		{name: "rate limited", status: http.StatusTooManyRequests},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client := newTestClient(t,
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, testAPIKey, r.Header.Get("AccessKey"))
					w.WriteHeader(tc.status)
					_, err := w.Write([]byte("request rejected"))
					assert.NoError(t, err)
				},
			)
			res, err := client.PullZonePublicIndexAllWithResponse(
				context.Background(), nil,
			)
			require.NoError(t, err)
			assert.Equal(t, tc.status, res.StatusCode())
			assert.Nil(t, res.JSON200)
			assert.Equal(t, "request rejected", string(res.Body))
		})
	}
}

func TestMalformedResponse(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte("{"))
		assert.NoError(t, err)
	})
	_, err := client.PullZonePublicIndexAllWithResponse(
		context.Background(), nil,
	)

	var syntaxError *json.SyntaxError

	require.ErrorAs(t, err, &syntaxError)
}

func TestCancelledRequest(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("cancelled request reached server")
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.PullZonePublicIndexAllWithResponse(ctx, nil)
	require.ErrorIs(t, err, context.Canceled)
}

func TestRawPullZonePagination(t *testing.T) {
	t.Parallel()

	const (
		pageSize   = 5
		totalItems = 13
	)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/pullzone", r.URL.Path)
		assert.Equal(t, testAPIKey, r.Header.Get("AccessKey"))
		assert.Equal(t, "5", r.URL.Query().Get("perPage"))

		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		assert.NoError(t, err)

		start := (page - 1) * pageSize
		items := make([]bunnynet.PullZoneModel, 0, pageSize)

		for index := start; index < min(start+pageSize, totalItems); index++ {
			items = append(items, bunnynet.PullZoneModel{
				Id: new(int64(index + 1)),
			})
		}

		w.Header().Set("Content-Type", "application/json")
		assert.NoError(t, json.NewEncoder(w).Encode(map[string]any{
			"Items":        items,
			"CurrentPage":  page,
			"TotalItems":   totalItems,
			"HasMoreItems": start+len(items) < totalItems,
		}))
	})

	identities := make([]int64, 0, totalItems)

	for _, page := range []int32{1, 2, 3, 4, 100} {
		res, err := client.PullZonePublicIndexAll(
			context.Background(),
			&bunnynet.PullZonePublicIndexAllParams{
				Page:    new(page),
				PerPage: new(int32(pageSize)),
			},
		)
		require.NoError(t, err)

		var envelope map[string]json.RawMessage

		err = json.NewDecoder(res.Body).Decode(&envelope)
		require.NoError(t, res.Body.Close())
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)

		var items []bunnynet.PullZoneModel

		require.NoError(t, json.Unmarshal(envelope["Items"], &items))

		remaining := totalItems - (int(page)-1)*pageSize
		require.Len(t, items, max(0, min(pageSize, remaining)))
		assert.JSONEq(
			t,
			strconv.Itoa(totalItems),
			string(envelope["TotalItems"]),
		)
		assert.JSONEq(
			t,
			strconv.Itoa(int(page)),
			string(envelope["CurrentPage"]),
		)
		assert.JSONEq(
			t,
			strconv.FormatBool(page < 3),
			string(envelope["HasMoreItems"]),
		)

		for _, item := range items {
			require.NotNil(t, item.Id)
			identities = append(identities, *item.Id)
		}
	}

	assert.Equal(
		t,
		[]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
		identities,
	)
}
