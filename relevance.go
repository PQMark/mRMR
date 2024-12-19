package mRMR

import (
	"math"
)

// Relevance computes the relevance of each feature with respect to the class and returns the scores as a slice.
func Relevance(data [][]float64, class []int, relevanceFunc func([]float64, []int) float64) []float64 {
	n := len(data[0])
	relevance := make([]float64, n)

	for i := 0; i < n; i++ {
		feature := getCol(data, i)
		mi := relevanceFunc(feature, class)
		relevance[i] = mi 
	}

	return relevance
}


// MutualInfo calculates the mutual information between two data slices.
func MutualInfo[T1, T2 Numeric](data1 []T1, data2 []T2) float64 {

	HA := shannonEntropy(data1)
	HB := shannonEntropy(data2)
	HAB := shannonJointEntropy(data1, data2)

	mi := HA + HB - HAB

	return mi
}

func shannonEntropy[T Numeric](sample []T) float64 {
	n := float64(len(sample))
	count := make(map[float64]int)
	sum := 0.0

	for _, val := range sample {
		count[float64(val)] ++
	}

	for _, val := range count {
		prob := float64(val) / n
		temp := prob * math.Log2(prob)
		sum += temp
	}

	return -sum 
}

func shannonJointEntropy[T1, T2 Numeric](data1 []T1, data2 []T2) float64 {
	if len(data1) != len(data2) {
		panic("Fail to calculate joint entropy: Unequal length of data")
	}

	n := float64(len(data1))
	count := make(map[[2]float64]int)
	sum := 0.0

	for i, val1 := range data1 {
		val2 := data2[i]

		data := [2]float64{float64(val1), float64(val2)}
		count[data]++
	}

	for _, val := range count {
		prob := float64(val) / n
		temp := prob * math.Log2(prob)
		sum += temp
	}

	return -sum
}

// FStatistic returns the f-statistic of feature and class. 
func FStatistic(feature []float64, class []int) float64 {
	bigN := float64(len(feature))

	normalized_ss := squaresOfSum(feature) / bigN
	ssbn := 0.0

	groups := groupByClass(feature, class)
	for _, g := range groups {
		ssbn += squaresOfSum(g) / float64(len(g))
	}
	ssbn -= normalized_ss

	mean := mean(feature)
	for i := range feature {
		feature[i] -= mean
	}

	sstotal := sumOfSquares(feature)

	sswn := sstotal - ssbn
	dfbn := float64(len(groups)) - 1  
	dfwn := bigN - float64(len(groups))
	
	msb := ssbn / dfbn
	msw := sswn / dfwn

	return msb / msw
}

func groupByClass(data []float64, class []int) [][]float64 {
	
	if len(data) != len(class) {
		panic("data and class slices must have the same length")
	}
	
	gmap := make(map[int][]float64)
	for i, cls := range class {
		gmap[cls] = append(gmap[cls], data[i])
	}

	grouped := make([][]float64, 0, len(gmap))
	for _, values := range gmap {
		grouped = append(grouped, values)
	}

	return grouped
}

// return (a + b + ...)^2
func squaresOfSum(data []float64) float64 {
	sum := 0.0

	for _, val := range data {
		sum += val 
	}

	return sum * sum
}

// return a^2 + b^2 + ...
func sumOfSquares(data []float64) float64 {
	sum := 0.0

	for _, val := range data {
		sum += val * val 
	} 

	return sum 
}
