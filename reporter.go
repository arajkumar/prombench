package prombench

import "io"

// Abstracts various type of reporters. e.g. Console, JSON, CSV..
type Reporter interface {
	Report(out io.Writer, s Summary) error
}
