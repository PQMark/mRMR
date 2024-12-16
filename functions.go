package mRMR

import (
	"fmt"
	"math"
)


func Discretization (data [][]float64, binSize int) [][]float64 {
	if len(data) == 0 || len(data[0]) == 0 {
		return nil 
	}

	r := len(data)
	c := len(data[0])

	min := make([]float64, c)
	max := make([]float64, c)

	for j := 0; j < c; j ++ {

		min[j] = data[0][j]
		max[j] = data[0][j]

		for i := 0; i < r; i++ {
			if data[i][j] > max[j] {
				max[j] = data[i][j]
			}
			if data[i][j] < min[j] {
				min[j] = data[i][j]
			}
		}
	}

	discreteData := make([][]float64, r)
	for i := range discreteData {
		discreteData[i] = make([]float64, c)
	}

	for j := 0; j < c; j++ {
		
		binWidth := (max[j] - min[j]) / float64(binSize)

		for i := 0; i < r; i++ {
			binIdx := int(math.Floor((data[i][j] - min[j]) / binWidth))

			if binIdx == binSize {
				binIdx--
			}

			discreteData[i][j] = float64(binIdx)
		}
	}

	return discreteData

}


func getCol(data [][]float64, i int) []float64 {
	col := make([]float64, len(data))

	for n, val := range data {
		col[n] = val[i]
	}

	return col
}


func selectByIndex[T any](data []T, idx []int) []T {
	var r []T 

	for _, index := range idx {
		if index < 0 || index >= len(data) {
			panic(fmt.Sprintf("index %d out of range", index))
		}

		r = append(r, data[index])
	}

	return r
}

func PairwiseDeduction[T Numeric](data1, data2 []T) []T {
	if len(data1) != len(data2) {
		panic("Fail to perform pairwise sum: Unequal length of data")
	}

	r := make([]T, len(data1))

	for i, val := range data1 {
		sum := val - data2[i]
		r[i] = sum
	}

	return r 
}

// return the index of the maximum value in a list 
func getMaxIndex[T Numeric] (data []T) int {
	a := data[0]
	idx := 0 

	for i, val := range data {
		if val > a {
			a = val 
			idx = i
		}
	}

	return idx
}


func Delete[T any] (data []T, idx int) []T {
	data = append(data[:idx], data[idx+1:]...)

	return data
}

func CheckIfAllNegative(data []float64) bool {
	for _,val := range data {
		if val > 0.0 && val > 0.001{
			return false
		}
	}
	return true
}