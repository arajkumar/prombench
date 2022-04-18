package prombench

type Query struct {
	PromQL string
	// Unix timestamp in milli.
	StartTime int64
	// Unix timestamp in milli.
	EndTime int64
	Step    int64
}
