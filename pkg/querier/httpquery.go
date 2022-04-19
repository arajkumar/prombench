package httpquery

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/arajkumar/prombench"
)

const (
	rfc3339Milli          = "2006-01-02T15:04:05.000Z07:00"
	queryRangeApiEndpoint = "/api/v1/query_range"
)

type HttpQuery struct {
	host    url.URL
	headers http.Header
	client  *http.Client
}

// Option controls the configuration of a Matcher.
type Option func(*HttpQuery) error

func New(host url.URL, opt ...Option) (HttpQuery, error) {
	h := HttpQuery{host: host}

	for _, f := range opt {
		if err := f(&h); err != nil {
			return h, err
		}
	}

	if h.client == nil {
		h.client = http.DefaultClient
	}

	return h, nil
}

// Returns http.Client
func (h HttpQuery) Client() *http.Client {
	return h.client
}

// WithClient sets the http.Client that the worker should use for requests.
// If not passed to NewHttpQuery, http.DefaultClient will be used.
func WithClient(c *http.Client) Option {
	return func(h *HttpQuery) error {
		h.client = c
		return nil
	}
}

// WithHeaders adds additional headers to http.Request.
func WithHeaders(headers http.Header) Option {
	return func(h *HttpQuery) error {
		h.headers = headers
		return nil
	}
}

func (h HttpQuery) NewHttpRequest(ctx context.Context, q prombench.Query) (*http.Request, error) {
	// Construct http request url.
	queryParam := url.Values{}
	queryParam.Set("query", q.PromQL)
	// TODO: Get rid of Query struct and replace it with url.Values type.
	// queryParam.Set("start", strconv.FormatFloat(float64(q.StartTime)/1000, 'f', 3, 64))
	// queryParam.Set("end", strconv.FormatFloat(float64(q.EndTime)/1000, 'f', 3, 64))
	queryParam.Set("start", time.UnixMilli(q.StartTime).UTC().Format(rfc3339Milli))
	queryParam.Set("end", time.UnixMilli(q.EndTime).UTC().Format(rfc3339Milli))
	queryParam.Set("step", strconv.FormatInt(q.Step, 10))
	host := h.host
	host.RawQuery = queryParam.Encode()
	host.Path = queryRangeApiEndpoint

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, host.String(), nil)
	if err != nil {
		return nil, err
	}
	// append additional headers.
	for hk, hv := range h.headers {
		r.Header.Set(hk, hv[0])
	}
	return r, nil
}
