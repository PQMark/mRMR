package main

import(
	"src/mRMR"
	"math"
    "testing"
)

func TestPearsonCorrelation(t *testing.T) {
	const epsilon = 1e-6

	testCases := []struct {
		data1  []float64
        data2    []float64
        expected float64
	}{
		{
			data1:  []float64{1,2,3,4,5},
			data2:    []float64{2,4,6,8,10},
			expected: 1.0,	
		},
		{
			data1: []float64{1,2,3,4,5},
			data2: []float64{2,2,2,2,2},
			expected: 0.0,
		},
		{
			data1: []float64{0,0,0,1,1},
			data2: []float64{1,1,1,0,0},
			expected: 1.0,  // -1.0 
		},
		{
			data1: []float64{1,2,3,4,5},
			data2: []float64{5,6,7,8,7},
			expected: 0.83205,
		},
	}

	for _, tt := range testCases {
		result := mRMR.PearsonCorrelation(tt.data1, tt.data2)

		if math.Abs(result-tt.expected) > epsilon {
			t.Errorf("Expected %v, got %v", tt.expected, result)
		}
	}

}