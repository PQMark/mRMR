package main

import (
	//"fmt"
	"fmt"
	"math"
)

type Numeric interface{
	int | int8 | int16 | int32 | int64 | float32 | float64
}

// for classification only now
// class is mapped to 0-??
// without extremely small value correction
func mRMR_Discrete(dataRaw [][]float64, class []int, binSize, maxFeatures int) ([]int, []float64, map[[2]int]float64){
	
	data := Discretization(dataRaw, binSize)

	relevanceAll := RelevanceMI(data, class)

	if maxFeatures > len(data[0]) {
		maxFeatures = len(data[0])
	}

	featuresToConsider := make([]int, 0, len(data[0]))	
	selectedFeatures := make([]int, 0, len(data[0]))
	redundancyMap := make(map[[2]int]float64)
	
	// Discard if relevance MI = 0
	for i, val := range relevanceAll {
		if val > 0.0 {
			featuresToConsider = append(featuresToConsider, i)
		}
	}

	for i := 0; i < maxFeatures; i++ {

		relevance := selectByIndex(relevanceAll, featuresToConsider)
		redundancy := make([]float64, len(featuresToConsider))

		if i == 0 {
			// Initialize 
			redundancy = make([]float64, len(featuresToConsider))   // all 0 
		} else {
			// Calculate redundancy 
			lastSelectedF := selectedFeatures[len(selectedFeatures)-1]

			// Update redundancyMap 
			redundancyMap = RedundancyMIUpdate(data, featuresToConsider, lastSelectedF, redundancyMap)

			// Calculate redundancy for each feature candidate 
			for i, val1 := range featuresToConsider {
				miSum := 0.0
				for _, val2 := range selectedFeatures {
					key := [2]int{val2, val1}
					miSum += redundancyMap[key]
				}

				redundancy[i] = miSum / float64(len(selectedFeatures))
			}

		}

		fmt.Println("relevance:", relevance)
		fmt.Println("Redundancy:", redundancy)

		score := PairwiseDeduction(relevance, redundancy)
		fmt.Println(score)

		if CheckIfAllNegative(score) {
			break
		}

		idx := getMaxIndex(score)
		feature := featuresToConsider[idx]
		selectedFeatures = append(selectedFeatures, feature)
		featuresToConsider = Delete(featuresToConsider, idx)

	}             

	return selectedFeatures, relevanceAll, redundancyMap
}


func Discretization(data [][]float64, binSize int) [][]int {
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

	discreteData := make([][]int, r)
	for i := range discreteData {
		discreteData[i] = make([]int, c)
	}

	for j := 0; j < c; j++ {
		
		binWidth := (max[j] - min[j]) / float64(binSize)

		for i := 0; i < r; i++ {
			binIdx := int(math.Floor((data[i][j] - min[j]) / binWidth))

			if binIdx == binSize {
				binIdx--
			}

			discreteData[i][j] = binIdx
		}
	}

	return discreteData

}

func RelevanceMI(data [][]int, class []int) []float64 {
	n := len(data[0])
	MI := make([]float64, n)

	for i := 0; i < n; i++ {
		feature := getCol(data, i)
		mi := MutualInfo(feature, class)
		MI[i] = mi 
	}

	return MI
}

// map [selected, not selected]MI
func RedundancyMIUpdate(data [][]int, featureToConsider []int, target int, redundancyMap map[[2]int]float64) map[[2]int]float64 {

	data2 := getCol(data, target)

	for _, idx := range featureToConsider {
		data1 := getCol(data, idx)
		mi := MutualInfo(data1, data2)
		redundancyMap[[2]int{target, idx}] = mi
	}

	return redundancyMap
}

func MutualInfo(data1 []int, data2 []int) float64 {

	HA := ShannonEntropy(data1)
	HB := ShannonEntropy(data2)
	HAB := ShannonJointEntropy(data1, data2)

	mi := HA + HB - HAB

	return mi
}

func ShannonEntropy(sample []int) float64 {
	n := float64(len(sample))
	count := make(map[int]int)
	sum := 0.0

	for _, val := range sample {
		count[val] ++
	}

	for _, val := range count {
		prob := float64(val) / n
		temp := prob * math.Log2(prob)
		sum += temp
	}

	return -sum 
}

func ShannonJointEntropy(data1, data2 []int) float64 {
	if len(data1) != len(data2) {
		panic("Fail to calculate joint entropy: Unequal length of data")
	}

	n := float64(len(data1))
	count := make(map[[2]int]int)
	sum := 0.0

	for i, val1 := range data1 {
		val2 := data2[i]

		data := [2]int{val1, val2}
		count[data]++
	}

	for _, val := range count {
		prob := float64(val) / n
		temp := prob * math.Log2(prob)
		sum += temp
	}

	return -sum
}


func getCol(data [][]int, i int) []int {
	col := make([]int, len(data))

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