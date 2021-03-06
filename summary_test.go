package prombench

import (
	"fmt"
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
				NumOfQueries:   1,
				TotalDuration:  1 * time.Second,
				Min:            1 * time.Second,
				Max:            1 * time.Second,
				Avg:            1 * time.Second,
				Median:         1 * time.Second,
				StatusCodeDist: map[int]int{0: 1},
				ErrDist:        map[string]int{},
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
				NumOfQueries:   2,
				TotalDuration:  3 * time.Second,
				Min:            1 * time.Second,
				Max:            2 * time.Second,
				Avg:            3 * time.Second / 2,
				Median:         3 * time.Second / 2,
				StatusCodeDist: map[int]int{0: 2},
				ErrDist:        map[string]int{},
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
				NumOfQueries:   3,
				TotalDuration:  8 * time.Second,
				Min:            1 * time.Second,
				Max:            5 * time.Second,
				Avg:            8 * time.Second / 3,
				Median:         2 * time.Second,
				StatusCodeDist: map[int]int{0: 3},
				ErrDist:        map[string]int{},
			},
		},
		{
			description: "various status code",
			in: Report{
				{
					Duration:   1 * time.Second,
					StatusCode: 202,
				},
				{
					Duration:   1 * time.Second,
					StatusCode: 200,
				},
				{
					Duration:   1 * time.Second,
					StatusCode: 301,
				},
				{
					Duration:   1 * time.Second,
					StatusCode: 401,
				},
				{
					Duration:   1 * time.Second,
					StatusCode: 429,
				},
				{
					Duration:   1 * time.Second,
					StatusCode: 500,
				},
			},
			expected: Summary{
				NumOfQueries:   6,
				TotalDuration:  6 * time.Second,
				Min:            1 * time.Second,
				Max:            1 * time.Second,
				Avg:            1 * time.Second,
				Median:         1 * time.Second,
				StatusCodeDist: map[int]int{200: 2, 300: 1, 400: 2, 500: 1},
				ErrDist:        map[string]int{},
			},
		},
		{
			description: "single err",
			in: Report{
				{
					Duration: 1 * time.Second,
					Error:    fmt.Errorf("error test"),
				},
				{
					Duration: 1 * time.Second,
				},
			},
			expected: Summary{
				NumOfQueries:   2,
				TotalDuration:  2 * time.Second,
				Min:            1 * time.Second,
				Max:            1 * time.Second,
				Avg:            1 * time.Second,
				Median:         1 * time.Second,
				StatusCodeDist: map[int]int{0: 2},
				ErrDist:        map[string]int{"error test": 1},
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
