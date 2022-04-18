package promqlworker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"

	"github.com/arajkumar/prombench"
	"gotest.tools/v3/assert"
)

func TestQuery(t *testing.T) {
	tt := []struct {
		description string
		query       prombench.Query
		expected    string
	}{
		{
			description: "empty Query",
			query:       prombench.Query{},
			expected:    "http://foo:9090/api/v1/query_range?end=1970-01-01T00%3A00%3A00.000Z&query=&start=1970-01-01T00%3A00%3A00.000Z&step=0",
		},
		{
			description: "simple Query",
			query: prombench.Query{
				PromQL:    "hello{}",
				StartTime: 100,
				EndTime:   200,
				Step:      10,
			},
			expected: "http://foo:9090/api/v1/query_range?end=1970-01-01T00%3A00%3A00.200Z&query=hello%7B%7D&start=1970-01-01T00%3A00%3A00.100Z&step=10",
		},
	}

	ctx := context.Background()
	host, err := url.Parse("http://foo:9090")
	if err != nil {
		t.Fatalf("Unable to parse url %s", err)
	}
	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			w, err := NewPromQLWorker(*host)
			if err != nil {
				t.Fatalf("Unable to construct New Worker %s", err)
				return
			}
			r, err := w.(promqlWorker).newHttpPromQuery(ctx, tc.query)
			if err != nil {
				t.Fatalf("Unable to construct http.Request %s", err)
				return
			}
			assert.DeepEqual(t, tc.expected, r.URL.String())
		})
	}
}

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
	w, _ := NewPromQLWorker(*url, WithClient(server.Client()))
	w.Run(ctx, inC)
	if count != int64(len(queries)) {
		t.Errorf("Expected to work %v times, found %v", len(queries), count)
	}
}
