package csvparser

import (
	"io"
	"strings"
	"testing"

	"github.com/arajkumar/prombench"
	"gotest.tools/v3/assert"
)

func TestCSVParser(t *testing.T) {
	tt := []struct {
		in          string
		concurrency int
		expected    []prombench.Query
		err         string
	}{
		{
			// empty with concurrency 1
			in:       "",
			expected: []prombench.Query{},
		},
		{
			// empty with concurrency 10
			concurrency: 10,
			in:          "",
			expected:    []prombench.Query{},
		},
		{
			// perfect line
			concurrency: 10,
			in:          `demo_cpu_usage_seconds_total{mode="idle"}|100|200|50`,
			expected: []prombench.Query{
				{
					PromQL:    `demo_cpu_usage_seconds_total{mode="idle"}`,
					StartTime: 100,
					EndTime:   200,
					Step:      50,
				},
			},
		},
		{
			// leading and trailing spaces
			concurrency: 10,
			in: `
demo_cpu_usage_seconds_total{mode="idle"} | 100 |200|50
`,
			expected: []prombench.Query{
				{
					PromQL:    `demo_cpu_usage_seconds_total{mode="idle"}`,
					StartTime: 100,
					EndTime:   200,
					Step:      50,
				},
			},
		},
		{
			// leading and trailing spaces
			concurrency: 10,
			in: `
demo_cpu_usage_seconds_total{mode="idle"} | 100 |200|50
A| 10 |20|5
`,
			expected: []prombench.Query{
				{
					PromQL:    `demo_cpu_usage_seconds_total{mode="idle"}`,
					StartTime: 100,
					EndTime:   200,
					Step:      50,
				},
				{
					PromQL:    `A`,
					StartTime: 10,
					EndTime:   20,
					Step:      5,
				},
			},
		},
	}

	for _, tc := range tt {
		csv, err := NewCSVParser(strings.NewReader(tc.in), WithConcurrency(tc.concurrency))
		if err != nil {
			t.Errorf("NewCSVParser failed %s", err)
		}
		err = csv.Parse()
		if err != nil && err != io.EOF {
			t.Errorf("Parse failed %s", err)
		}

		actual := []prombench.Query{}
		for q := range csv.Query() {
			actual = append(actual, q)
		}
		assert.DeepEqual(t, tc.expected, actual)
	}
}
