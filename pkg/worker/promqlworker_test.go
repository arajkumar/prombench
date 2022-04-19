package promqlworker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"

	"github.com/arajkumar/prombench"
	httpquery "github.com/arajkumar/prombench/pkg/querier"
)

func TestPromqlWorker(t *testing.T) {
	ctx := context.Background()
	var count int64
	handler := func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&count, int64(1))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	queries := []prombench.Query{
		{
			PromQL:    "foo",
			StartTime: 100,
			EndTime:   200,
			Step:      10,
		},
		{
			PromQL:    "bar",
			StartTime: 100,
			EndTime:   200,
			Step:      10,
		},
	}
	inC := make(chan prombench.Query, len(queries))
	for _, q := range queries {
		inC <- q
	}
	close(inC)

	url, _ := url.Parse(server.URL)
	h, _ := httpquery.New(*url, httpquery.WithClient(server.Client()))
	w, _ := New(h, 10)
	w.Run(ctx, inC)
	if count != int64(len(queries)) {
		t.Errorf("Expected to work %v times, found %v", len(queries), count)
	}
}
