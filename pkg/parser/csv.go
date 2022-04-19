package csvparser

import (
	"context"
	"encoding/csv"
	"io"
	"strconv"
	"strings"

	"github.com/arajkumar/prombench"
)

// Option controls the configuration of a Matcher.
type Option func(*csvParser) error

type csvParser struct {
	out       chan prombench.Query
	csvReader *csv.Reader
	chSize    int
}

func New(r io.Reader, opt ...Option) (prombench.Parser, error) {
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

	// defaults to sane channel capacity.
	if c.chSize < 1 {
		c.chSize = 100
	}

	c.out = make(chan prombench.Query, c.chSize)

	return c, nil
}

// WithChannelSize sets the channel capacity.
//
// If not passed to NewCSVParser, sane default(100) will be used.
func WithChannelSize(chSize int) Option {
	return func(c *csvParser) error {
		c.chSize = chSize
		return nil
	}
}

// Implements prombench.Parser interface.
func (c csvParser) Queries() <-chan prombench.Query {
	return c.out
}

// Implements prombench.Parser interface.
func (c csvParser) Parse(ctx context.Context) error {
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
		select {
		case c.out <- prombench.Query{
			PromQL:    strings.TrimSpace(record[0]),
			StartTime: startTime,
			EndTime:   endTime,
			Step:      step,
		}:
			continue
		case <-ctx.Done():
			break
		}
	}
	return nil
}
