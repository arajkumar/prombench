package csvparser

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/arajkumar/prombench"
	"gotest.tools/v3/assert"
)

func TestCSVParser(t *testing.T) {
	tt := []struct {
		description string
		in          string
		chSize      int
		expected    []prombench.Query
		err         string
	}{
		{
			description: "empty with chSize 1",
			in:          "",
			expected:    []prombench.Query{},
		},
		{
			description: "empty with chSize 10",
			chSize:      10,
			in:          "",
			expected:    []prombench.Query{},
		},
		{
			description: "perfect line",
			chSize:      10,
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
			description: "leading and trailing spaces 1",
			chSize:      10,
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
			description: "leading and trailing spaces 2",
			chSize:      10,
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

	ctx := context.Background()
	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			csv, err := NewCSVParser(strings.NewReader(tc.in), WithChannelSize(tc.chSize))
			if err != nil {
				t.Errorf("NewCSVParser failed %s", err)
			}
			err = csv.Parse(ctx)
			if err != nil && err != io.EOF {
				t.Errorf("Parse failed %s", err)
			}

			actual := []prombench.Query{}
			for q := range csv.Queries() {
				actual = append(actual, q)
			}
			assert.DeepEqual(t, tc.expected, actual)
		})
	}
}
