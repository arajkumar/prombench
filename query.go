package prombench

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

type Query struct {
	PromQL    string
	StartTime int64
	EndTime   int64
	Step      int64
}

type QueryChannel chan Query

func (q Query) NewHttpPromQuery(ctx context.Context, host *url.URL) (*http.Request, error) {
	// Construct http request url.
	queryParam := url.Values{}
	queryParam.Set("query", q.PromQL)
	// TODO: Get rid of Query struct and replace it with url.Values type.
	queryParam.Set("start", strconv.FormatInt(q.StartTime, 10))
	queryParam.Set("end", strconv.FormatInt(q.EndTime, 10))
	queryParam.Set("step", strconv.FormatInt(q.Step, 10))
	h := host
	h.RawQuery = queryParam.Encode()
	h.Path = "/api/v1/query_range"

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, h.String(), nil)
	if err != nil {
		return nil, err
	}
	return r, nil
}
