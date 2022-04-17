package prombench

import (
	"context"
	"testing"

	"net/url"

	"gotest.tools/v3/assert"
)

func TestQuery(t *testing.T) {
	tt := []struct {
		query    Query
		expected string
	}{
		{
			query:    Query{},
			expected: "http://foo:9090?end=0&query=&start=0&step=0",
		},
		{
			query: Query{
				PromQL:    "hello{}",
				StartTime: 100,
				EndTime:   200,
				Step:      10,
			},
			expected: "http://foo:9090?end=200&query=hello%7B%7D&start=100&step=10",
		},
	}

	ctx := context.Background()
	host, err := url.Parse("http://foo:9090")
	if err != nil {
		t.Fatalf("Unable to parse url %s", err)
	}
	for _, tc := range tt {
		r, err := tc.query.NewHttpPromQuery(ctx, host)
		if err != nil {
			t.Fatalf("Unable to construct http.Request %s", err)
			continue
		}
		assert.DeepEqual(t, tc.expected, r.URL.String())
	}
}
