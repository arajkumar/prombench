package httpquery

import (
	"context"
	"net/url"
	"testing"

	"github.com/arajkumar/prombench"
	"gotest.tools/v3/assert"
)

func TestHttpQuery(t *testing.T) {
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
			h, err := New(*host)
			if err != nil {
				t.Fatalf("Unable to construct New Worker %s", err)
				return
			}
			r, err := h.NewHttpRequest(ctx, tc.query)
			if err != nil {
				t.Fatalf("Unable to construct http.Request %s", err)
				return
			}
			assert.DeepEqual(t, tc.expected, r.URL.String())
		})
	}
}
