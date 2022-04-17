package csvparser

import (
	"encoding/csv"
	"io"
	"strconv"
	"strings"

	"github.com/arajkumar/prombench"
)

// Option controls the configuration of a Matcher.
type Option func(*csvParser) error

type csvParser struct {
	out         chan prombench.Query
	csvReader   *csv.Reader
	concurrency int
}

func NewCSVParser(r io.Reader, opt ...Option) (prombench.Parser, error) {
	csvReader := csv.NewReader(r)
	// Data file would contain | as separator instead of default(,).
	csvReader.Comma = '|'
	csvReader.LazyQuotes = true
	c := csvParser{
		csvReader: csvReader,
	}

	for _, f := range opt {
		if err := f(&c); err != nil {
			return nil, err
		}
	}

	// defaults to sane concurrency limit.
	if c.concurrency < 1 {
		c.concurrency = 1
	}

	c.out = make(prombench.QueryChannel, c.concurrency)

	return c, nil
}

// WithConcurrency sets the concurrency limit for the network calls.
//
// If not passed to NewPromQLWorker, sane default(1) will be used.
func WithConcurrency(concurrency int) Option {
	return func(c *csvParser) error {
		c.concurrency = concurrency
		return nil
	}
}

// Implements prombench.Parser interface.
func (c csvParser) Query() prombench.QueryChannel {
	return c.out
}

// Implements prombench.Parser interface.
func (c csvParser) Parse() error {
	defer close(c.out)
	for {
		record, err := c.csvReader.Read()
		if err != nil {
			return err
		}
		startTime, err := strconv.ParseInt(strings.TrimSpace(record[1]), 10, 64)
		if err != nil {
			return err
		}
		endTime, err := strconv.ParseInt(strings.TrimSpace(record[2]), 10, 64)
		if err != nil {
			return err
		}
		step, err := strconv.ParseInt(strings.TrimSpace(record[3]), 10, 64)
		if err != nil {
			return err
		}
		c.out <- prombench.Query{
			PromQL:    strings.TrimSpace(record[0]),
			StartTime: startTime,
			EndTime:   endTime,
			Step:      step,
		}
	}
	return nil
}
