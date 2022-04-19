package promqlworker

import (
	"context"
	"io"
	"io/ioutil"
	"sync"
	"time"

	"github.com/arajkumar/prombench"
	httpquery "github.com/arajkumar/prombench/pkg/querier"
)

// Option controls the configuration of a Matcher.
type Option func(*promqlWorker) error

type promqlWorker struct {
	concurrency int
	httpQuery   httpquery.HttpQuery
}

func New(httpQuery httpquery.HttpQuery, concurrency int) (prombench.Worker, error) {
	// Serial worker!
	if concurrency < 1 {
		concurrency = 1
	}

	return promqlWorker{
		concurrency: concurrency,
		httpQuery:   httpQuery,
	}, nil
}

func (w promqlWorker) runWorker(ctx context.Context, q prombench.Query) prombench.Stat {
	r, err := w.httpQuery.NewHttpRequest(ctx, q)
	if err != nil {
		return prombench.Stat{Error: err}
	}
	// TODO: should we exclude other latencies like DNS dial..?
	start := time.Now()
	resp, err := w.httpQuery.Client().Do(r)
	if err == nil {
		// Read and ignore reponse body!
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}
	return prombench.Stat{
		Duration: time.Since(start),
		Error:    err,
	}
}

func (w promqlWorker) run(ctx context.Context, queries <-chan prombench.Query, stat chan<- prombench.Stat) {
	for q := range queries {
		stat <- w.runWorker(ctx, q)
	}
}

// Implements Worker interface.
func (w promqlWorker) Run(ctx context.Context, queries <-chan prombench.Query) prombench.Report {
	statC := make(chan prombench.Stat, cap(queries))

	var wg sync.WaitGroup
	wg.Add(w.concurrency)

	for i := 0; i < w.concurrency; i++ {
		go func() {
			w.run(ctx, queries, statC)
			wg.Done()
		}()
	}

	// Wait for all workers to finish and then close statC
	go func() {
		wg.Wait()
		close(statC)
	}()

	report := prombench.Report{}
	for stat := range statC {
		report = append(report, stat)
	}
	return report
}
