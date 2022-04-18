package promqlworker

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/arajkumar/prombench"
	"golang.org/x/sync/errgroup"
)

// Option controls the configuration of a Matcher.
type Option func(*promqlWorker) error

type promqlWorker struct {
	client      *http.Client
	concurrency int
	headers     http.Header
	host        url.URL
}

func NewPromQLWorker(host url.URL, opt ...Option) (prombench.Worker, error) {
	w := promqlWorker{host: host}

	for _, f := range opt {
		if err := f(&w); err != nil {
			return nil, err
		}
	}

	if w.client == nil {
		w.client = http.DefaultClient
	}
	// defaults to sane concurrency limit.
	if w.concurrency < 1 {
		w.concurrency = 1
	}

	return w, nil
}

// WithClient sets the http.Client that the worker should use for requests.
// If not passed to NewPromQLWorker, http.DefaultClient will be used.
func WithClient(c *http.Client) Option {
	return func(w *promqlWorker) error {
		w.client = c
		return nil
	}
}

// WithConcurrency sets the concurrency limit for the network calls.
//
// If not passed to NewPromQLWorker, sane default(1) will be used.
func WithConcurrency(concurrency int) Option {
	return func(w *promqlWorker) error {
		w.concurrency = concurrency
		return nil
	}
}

// WithHeaders adds additional headers to http.Request.
func WithHeaders(headers http.Header) Option {
	return func(w *promqlWorker) error {
		w.headers = headers
		return nil
	}
}

const rfc3339Milli = "2006-01-02T15:04:05.000Z07:00"

func (w promqlWorker) newHttpPromQuery(ctx context.Context, q prombench.Query) (*http.Request, error) {
	// Construct http request url.
	queryParam := url.Values{}
	queryParam.Set("query", q.PromQL)
	// TODO: Get rid of Query struct and replace it with url.Values type.
	// queryParam.Set("start", strconv.FormatFloat(float64(q.StartTime)/1000, 'f', 3, 64))
	// queryParam.Set("end", strconv.FormatFloat(float64(q.EndTime)/1000, 'f', 3, 64))
	queryParam.Set("start", time.UnixMilli(q.StartTime).UTC().Format(rfc3339Milli))
	queryParam.Set("end", time.UnixMilli(q.EndTime).UTC().Format(rfc3339Milli))
	queryParam.Set("step", strconv.FormatInt(q.Step, 10))
	h := w.host
	h.RawQuery = queryParam.Encode()
	h.Path = "/api/v1/query_range"

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, h.String(), nil)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (w promqlWorker) run(ctx context.Context, q prombench.Query) (prombench.Stat, error) {
	r, err := w.newHttpPromQuery(ctx, q)
	if err != nil {
		return prombench.Stat{}, err
	}
	// append additional headers.
	for hk, hv := range w.headers {
		r.Header.Set(hk, hv[0])
	}
	start := time.Now()
	resp, err := w.client.Do(r)
	if err == nil {
		// Read and ignore reponse body!
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	} else {
		return prombench.Stat{}, err
	}
	return prombench.Stat{
		Duration: time.Since(start),
	}, err
}

// Implements Worker interface.
func (w promqlWorker) Run(ctx context.Context, queries <-chan prombench.Query) (prombench.Report, error) {
	errC := make(chan error, 1)
	ctrlC := make(chan prombench.Stat, cap(queries))
	go func() {
		defer close(errC)
		defer close(ctrlC)
		var g errgroup.Group
		for q := range queries {
			func(q prombench.Query) {
				g.Go(func() error {
					stat, err := w.run(ctx, q)
					if err != nil {
						return err
					}
					ctrlC <- stat
					return nil
				})
			}(q)
		}
		if err := g.Wait(); err != nil {
			errC <- err
		}
	}()

	report := prombench.Report{}
	for stat := range ctrlC {
		report = append(report, stat)
	}
	err := <-errC // guaranteed to have an err or be closed.
	return report, err
}
