package main

import (
	"fmt"
	"src/mRMR"
	"testing"
)


func TestDiscretization(t *testing.T) {
	
	tests := []struct{
		data     [][]float64
		binSize  int
		expected [][]float64
	} {
		{
			data: [][]float64{
				{1.5, 2.3, 3.8},
				{2.0, 3.1, 4.2},
				{1.8, 2.7, 4.0},
			},
			binSize: 3,
			expected: [][]float64{
				{0, 0, 0}, 
				{2, 2, 2}, 
				{1, 1, 1}, 
			},
		},
		{
			data: [][]float64{
				{1.5, 2.3},
				{2.0, 3.1},
			},
			binSize: 2,
			expected: [][]float64{
				{0, 0},
				{1, 1},
			},
		},
		{
			data: [][]float64{
				{0.5, 1.0, 1.5},
				{1.6, 2.1, 2.6},
			},
			binSize: 2,
			expected: [][]float64{
				{0, 0, 0},
				{1, 1, 1},
			},
		},
		{
			data: [][]float64{
				{0.1, 0.5},  
				{10.2, 2.3}, 
				{5.5, 8.7},  
			},
			binSize: 4,
			expected: [][]float64{
				{0, 0}, 
				{3, 0},
				{2, 3}, 
			},
		},
	}
	
	for _, tt := range tests {
		discreteData := mRMR.Discretization(tt.data, tt.binSize)

		for i := range discreteData {
			for j := range discreteData[i] {
				if discreteData[i][j] != tt.expected[i][j] {
					fmt.Println(discreteData)
					t.Errorf("Discretization failed at row %d, column %d. Got %v, expected %v", i, j, discreteData[i][j], tt.expected[i][j])
				}
			}
		}
	}
}

