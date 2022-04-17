package prombench

// Abstracts a work.
type Work func() error

// Abstracts bunch of concurrently executable work.
type Works interface {
	Next() (Work, error)
}

// Execution statistics of works.
type Report struct {
	durtion []float64
}

type Worker interface {
	Run(works Works) (Report, error)
}
