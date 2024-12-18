package main

import (
	"fmt"
	"math/rand"
	"src/mRMR"
	"testing"
	"math"
)

// defual value for binSize 

func TestMRMR(t *testing.T) {
	data := GenerateData(2000)

	paras := mRMR.ParasmRMR{
		Data: data,
		Discretization: false,
		Calculation: "diff",
		Method: "nmi-nmi",
		BinSize: 300,
	}

	selectedFeatures, _, _ := paras.MRMR()

	fmt.Println(selectedFeatures)
}

func TestDiscretization(t *testing.T) {
	
	tests := []struct{
		data     [][]float64
		binSize  int
		expected [][]float64
		quantized [][]float64
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
			quantized: [][]float64{
				{1.583, 2.433, 3.867},
				{1.917, 2.967, 4.133},
				{1.75,  2.7,   4.0},
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
			quantized: [][]float64{
				{1.625, 2.5},
				{1.875, 2.9},
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
			quantized: [][]float64{
				{0.775, 1.275, 1.775}, 
   				{1.325, 1.825, 2.325}, 
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
			quantized: [][]float64{
				{1.3625, 1.525},  
				{8.9375, 1.525},  
				{6.4125, 7.675}, 
			},	
		},
	}
	
	for _, tt := range tests {
		discreteData, quantizedData := mRMR.Discretization(tt.data, tt.binSize)

		for i := range discreteData {
			for j := range discreteData[i] {
				if discreteData[i][j] != tt.expected[i][j] {
					fmt.Println(discreteData)
					t.Errorf("Discretization failed at row %d, column %d. Got %v, expected %v", i, j, discreteData[i][j], tt.expected[i][j])
				}

				if math.Abs(quantizedData[i][j] - tt.quantized[i][j]) > 0.001 {
					fmt.Println(quantizedData)
					t.Errorf("Quantization failed at row %d, column %d. Got %v, expected %v", i, j, quantizedData[i][j], tt.quantized[i][j])
				}

				
			}
		}
	}
}


func GenerateData(nSamples int) mRMR.DatamRMR {
	r := rand.New(rand.NewSource(66))

	X := make([][]float64, nSamples)
	class := make([]int, nSamples)

	for i := 0; i < nSamples; i++ {
		X[i] = make([]float64, 6)    // 6 features: 2 relevant, 2 redundant, 2 noisy

		feature1 := r.NormFloat64()*0.5 + 0.8
		feature2 := r.Float64()			// interction with feature1, but will not stand out 

		if math.Sin(feature1*math.Pi)+feature2 > 1.0 {
			class[i] = 1
		} else {
			class[i] = 0
		}

		X[i][0] = feature1
		X[i][1] = feature2

		// Insert redundant features
		X[i][2] = feature1 + 0.01*r.Float64()
		X[i][3] = feature2 + 0.01*r.Float64()

		// Insert noise features
		X[i][4] = r.ExpFloat64()
		X[i][5] = r.Float64() 
	}

	return mRMR.DatamRMR{X: X, Class: class}
}