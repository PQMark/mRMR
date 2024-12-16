package main

import(
	"src/mRMR"
	"math"
    "testing"
)

func TestFStatistic(t *testing.T) {

	const epsilon = 1e-6

	testCases := []struct {
		feature  []float64
        class    []int
        expected float64
	}{
		{
			feature:  []float64{2, 4, 6, 5, 5, 5, 7, 8, 9},
			class:    []int{1, 1, 1, 2, 2, 2, 3, 3, 3},
			expected: 7.79999999,	
		},
		{
			feature: []float64{1.0, 2.0, 3.0, 2.0, 4.0, 6.0, 7.0, 8.0, 9.0},
			class: []int{1, 1, 1, 0, 0, 0, 3, 3, 3},
			expected: 14.000001,
		},
	}

	for _, tt := range testCases {
		result := mRMR.FStatistic(tt.feature, tt.class)

		if math.Abs(result-tt.expected) > epsilon {
			t.Errorf("Expected %v, got %v", tt.expected, result)
		}
	}

}