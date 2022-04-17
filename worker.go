package prombench

import (
	"context"
)

type Worker interface {
	Run(ctx context.Context, p Parser) (Report, error)
}
