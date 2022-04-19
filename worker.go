package prombench

import (
	"context"
)

// TODO: Check whether github.com/rakyll/hey can be used
// to replace this worker logic after https://github.com/rakyll/hey/pull/149.
type Worker interface {
	Run(ctx context.Context, queries <-chan Query) Report
}
