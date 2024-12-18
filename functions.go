package mRMR

import (
	"fmt"
	"math"
)

func Discretization (data [][]float64, binSize int) ([][]float64, [][]float64) {
	if len(data) == 0 || len(data[0]) == 0 {
		return nil, nil 
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
	quantizedData := make([][]float64, r)

	for i := range discreteData {
		discreteData[i] = make([]float64, c)
		quantizedData[i] = make([]float64, c)
	}

	for j := 0; j < c; j++ {
		
		binWidth := (max[j] - min[j]) / float64(binSize)

		for i := 0; i < r; i++ {
			binIdx := int(math.Floor((data[i][j] - min[j]) / binWidth))

			if binIdx == binSize {
				binIdx--
			}

			// replaced by bin indices
			discreteData[i][j] = float64(binIdx)
			
			// replaced by midpoints of bins
			binMidpoint := min[j] + (float64(binIdx)+0.5)*binWidth
			quantizedData[i][j] = binMidpoint
		}
	}

	return discreteData, quantizedData

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

func PairwiseOperation(data1, data2 []float64, operation string) []float64 {
	if len(data1) != len(data2) {
		panic("Fail to perform pairwise sum: Unequal length of data")
	}

	const epsilon = 1e-8
	r := make([]float64, len(data1))

	for i, val := range data1 {
		switch operation {
		case "diff":
			r[i] = val - data2[i]
		case "quo":
			divisor := data2[i]
			if divisor == 0 {
                divisor = epsilon //avoid division by zero
            }

			r[i] = val / divisor
		default:
			panic("Invalid operation. Choose from 'diff' or 'quo'")
		}
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
		if val > 0.0{
			return false
		}
	}
	return true
}

func CheckIfAllSmallerOne(data []float64) bool {
	for _, val := range data {
		if val > 1.0 {
			return false
		}
	}

	return true
}

func MinMaxNormalization(data []float64) []float64 {
	
	min := 0.0
	max := 0.0
	new := make([]float64, len(data))

	for _, val := range data {
		if val > max {
			max = val
		}

		if val < min {
			min = val
		}
	}

	diff := max - min

	for i, val := range data {
		new[i] = (val - min) / diff
	}

	return new
}

// get the quantization level
func QuantizationLevel(data [][]float64, threshold float64) int {
	level := 2 

	numFeatures := len(data[0])

	for {
		maxError := 0.0

		_, quantizedData := Discretization(data, level)

		for i := 0; i < numFeatures; i++ {
			quantizedFeature := getCol(quantizedData, i)
			originalFeature := getCol(data, i)

			QError := QuantizationError(quantizedFeature, originalFeature)

			if QError > maxError {
				maxError = QError
			}

		}

		if maxError <= threshold {
			return level
		}

		level++
	}
}

// get the quantization error
func QuantizationError(quantizedData, originalData []float64) float64 {
	n := float64(len(originalData))
	err := 0.0

	for i := range originalData {
		err += math.Pow(quantizedData[i] - originalData[i], 2)
	}

	return err / n
}

func scaling(data []float64, factor float64) []float64 {
	for i := range data {
		data[i] /= factor
	}

	return data
}

func uniqueClass(data []int) int {
	m := make(map[int]int)

	for _, val := range data {
		m[val] ++
	}

	return len(m)
}

func GetFeatures(features []string, indices []int) []string {
	selectedFeatures := make([]string, len(indices))

	for i, idx := range indices {
		selectedFeatures[i] = features[idx]
	}

	return selectedFeatures
}