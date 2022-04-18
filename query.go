package prombench

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Query struct {
	PromQL string
	// Unix timestamp in milli.
	StartTime int64
	// Unix timestamp in milli.
	EndTime int64
	Step    int64
}

const rfc3339Milli = "2006-01-02T15:04:05.000Z07:00"

func (q Query) NewHttpPromQuery(ctx context.Context, host url.URL) (*http.Request, error) {
	// Construct http request url.
	queryParam := url.Values{}
	queryParam.Set("query", q.PromQL)
	// TODO: Get rid of Query struct and replace it with url.Values type.
	// queryParam.Set("start", strconv.FormatFloat(float64(q.StartTime)/1000, 'f', 3, 64))
	// queryParam.Set("end", strconv.FormatFloat(float64(q.EndTime)/1000, 'f', 3, 64))
	queryParam.Set("start", time.UnixMilli(q.StartTime).UTC().Format(rfc3339Milli))
	queryParam.Set("end", time.UnixMilli(q.EndTime).UTC().Format(rfc3339Milli))
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
