package prombench

import (
	"context"
	"net/http"
)

// Converts prombench.Query into http.Request.
type HttpPromQuerier interface {
	HttpPromQuery(ctx context.Context, query Query) http.Request
}
