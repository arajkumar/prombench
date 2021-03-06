package prombench

import "context"

// Abstracts bunch of concurrently executable work items.
type Parser interface {
	// Runs a parse loop in a blocking mode.
	Parse(ctx context.Context) error
	// Thread safe iterator like construct returns Work item one by one.
	// Return Go channel which ferries Query object.
	Queries() <-chan Query
}
