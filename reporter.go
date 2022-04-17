package prombench

import "time"

// Execution statistics of works.
type Report struct {
	Duration []time.Duration
}

// Abstracts various type of reporters. e.g. Console, JSON, CSV..
type Reporter interface {
	Report(report Report) error
}
