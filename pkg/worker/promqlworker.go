package promqlworker

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/arajkumar/prombench"
	"golang.org/x/sync/errgroup"
)

// Option controls the configuration of a Matcher.
type Option func(*promqlWorker) error

type promqlWorker struct {
	client      *http.Client
	concurrency int
}

func NewPromQLWorker(opt ...Option) (prombench.Worker, error) {
	w := promqlWorker{}

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

func (w promqlWorker) run(ctx context.Context, host *url.URL, q prombench.Query) (time.Duration, error) {
	r, err := q.NewHttpPromQuery(ctx, host)
	if err != nil {
		return 0, err
	}
	start := time.Now()
	resp, err := w.client.Do(r)
	if err == nil {
		// Read and ignore reponse body!
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	} else {
		return 0, err
	}
	return time.Since(start), err
}

// Implements Worker interface.
func (w promqlWorker) Run(ctx context.Context, host *url.URL, queries prombench.QueryChannel) (prombench.Report, error) {
	errC := make(chan error, 1)
	ctrlC := make(chan time.Duration, 1)
	go func() {
		defer close(errC)
		defer close(ctrlC)
		var g errgroup.Group
		for q := range queries {
			g.Go(func() error {
				dur, err := w.run(ctx, host, q)
				if err != nil {
					return err
				}
				ctrlC <- dur
				return nil
			})
		}
		if err := g.Wait(); err != nil {
			errC <- err
		}
	}()

	duration := []time.Duration{}
	for d := range ctrlC {
		duration = append(duration, d)
	}
	err := <-errC // guaranteed to have an err or be closed.
	return prombench.Report{Duration: duration}, err
}
