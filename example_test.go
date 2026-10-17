package bunnynet_test

import (
	"context"
	"net/http"
	"os"

	"github.com/psyb0t/go-bunnynet"
)

//nolint:testableexamples // Needs real credentials, compile check only.
func Example() {
	auth := bunnynet.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set("AccessKey", os.Getenv("BUNNYNET_API_KEY"))

			return nil
		},
	)

	client, err := bunnynet.NewClientWithResponses(
		"https://api.bunny.net",
		auth,
	)
	if err != nil {
		panic(err)
	}

	res, err := client.PullZonePublicIndexAllWithResponse(
		context.Background(),
		nil,
	)
	if err != nil {
		panic(err)
	}

	_ = res.JSON200
}
