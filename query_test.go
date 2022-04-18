package prombench

import (
	"context"
	"testing"

	"net/url"

	"gotest.tools/v3/assert"
)

func TestQuery(t *testing.T) {
	tt := []struct {
		description string
		query       Query
		expected    string
	}{
		{
			description: "empty Query",
			query:       Query{},
			expected:    "http://foo:9090/api/v1/query_range?end=1970-01-01T05%3A30%3A00.000%2B05%3A30&query=&start=1970-01-01T05%3A30%3A00.000%2B05%3A30&step=0",
		},
		{
			description: "simple Query",
			query: Query{
				PromQL:    "hello{}",
				StartTime: 100,
				EndTime:   200,
				Step:      10,
			},
			expected: "http://foo:9090/api/v1/query_range?end=1970-01-01T05%3A30%3A00.200%2B05%3A30&query=hello%7B%7D&start=1970-01-01T05%3A30%3A00.100%2B05%3A30&step=10",
		},
	}

	ctx := context.Background()
	host, err := url.Parse("http://foo:9090")
	if err != nil {
		t.Fatalf("Unable to parse url %s", err)
	}
	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			r, err := tc.query.NewHttpPromQuery(ctx, *host)
			if err != nil {
				t.Fatalf("Unable to construct http.Request %s", err)
				return
			}
			assert.DeepEqual(t, tc.expected, r.URL.String())
		})
	}
}
