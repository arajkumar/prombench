package prombench

import (
	"sort"
	"time"
)

// Execution statistics of works.
type Stat struct {
	Duration time.Duration
}

type Report []Stat

type Summary struct {
	NumOfQueries  int
	TotalDuration time.Duration
	Min           time.Duration
	Max           time.Duration
	Avg           time.Duration
	Median        time.Duration
}

func (a Report) Len() int           { return len(a) }
func (a Report) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Report) Less(i, j int) bool { return a[i].Duration < a[j].Duration }

func (r Report) median() time.Duration {
	l := len(r)
	if l%2 == 0 {
		return (r[l/2].Duration + r[l/2-1].Duration) / 2
	}
	return r[l/2].Duration
}

func (r Report) ToSummary() Summary {
	sort.Sort(r)
	var totalDuration time.Duration
	for _, d := range r {
		totalDuration += d.Duration
	}
	l := len(r)
	s := Summary{
		NumOfQueries:  l,
		TotalDuration: totalDuration,
		Min:           r[0].Duration,
		Max:           r[l-1].Duration,
		Median:        r.median(),
		Avg:           totalDuration / time.Duration(l),
	}
	return s
}
