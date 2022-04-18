package prombench

// Abstracts various type of reporters. e.g. Console, JSON, CSV..
type Reporter interface {
	Report(s Summary) error
}
