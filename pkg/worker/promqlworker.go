package promqlworker

import (
	"context"
	"net/http"

	"github.com/arajkumar/prombench"
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

// Implements Worker interface.
func (w promqlWorker) Run(ctx context.Context, p prombench.Parser) (prombench.Report, error) {
	return prombench.Report{}, nil
}
