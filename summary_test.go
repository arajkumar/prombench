package prombench

import (
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func TestSummary(t *testing.T) {
	tt := []struct {
		description string
		in          Report
		expected    Summary
	}{
		{
			description: "single stat",
			in: Report{
				{
					Duration: 1 * time.Second,
				},
			},
			expected: Summary{
				NumOfQueries:  1,
				TotalDuration: 1 * time.Second,
				Min:           1 * time.Second,
				Max:           1 * time.Second,
				Avg:           1 * time.Second,
				Median:        1 * time.Second,
			},
		},
		{
			description: "even stats",
			in: Report{
				{
					Duration: 1 * time.Second,
				},
				{
					Duration: 2 * time.Second,
				},
			},
			expected: Summary{
				NumOfQueries:  2,
				TotalDuration: 3 * time.Second,
				Min:           1 * time.Second,
				Max:           2 * time.Second,
				Avg:           3 * time.Second / 2,
				Median:        3 * time.Second / 2,
			},
		},
		{
			description: "odd stats",
			in: Report{
				{
					Duration: 1 * time.Second,
				},
				{
					Duration: 2 * time.Second,
				},
				{
					Duration: 5 * time.Second,
				},
			},
			expected: Summary{
				NumOfQueries:  3,
				TotalDuration: 8 * time.Second,
				Min:           1 * time.Second,
				Max:           5 * time.Second,
				Avg:           8 * time.Second / 3,
				Median:        2 * time.Second,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			actual := tc.in.ToSummary()
			assert.DeepEqual(t, tc.expected, actual)
		})
	}
}
