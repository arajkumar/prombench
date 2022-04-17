package prombench

// Execution statistics of works.
type Report struct {
	durtion []float64
}

// Abstracts various type of reporters. e.g. Console, JSON, CSV..
type Reporter interface {
	Report(report Report) error
}
