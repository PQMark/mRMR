package mRMR

import (
	"math"
)

// 
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


// MI 
func MutualInfo[T1, T2 Numeric](data1 []T1, data2 []T2) float64 {

	HA := ShannonEntropy(data1)
	HB := ShannonEntropy(data2)
	HAB := ShannonJointEntropy(data1, data2)

	mi := HA + HB - HAB

	return mi
}

func ShannonEntropy[T Numeric](sample []T) float64 {
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

func ShannonJointEntropy[T1, T2 Numeric](data1 []T1, data2 []T2) float64 {
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

// F statistics
func FStatistic(feature []float64, class []int) float64 {
	bigN := float64(len(feature))

	normalized_ss := SquaresofSum(feature) / bigN
	ssbn := 0.0

	groups := GroupByClass(feature, class)
	for _, g := range groups {
		ssbn += SquaresofSum(g) / float64(len(g))
	}
	ssbn -= normalized_ss

	mean := Mean(feature)
	for i := range feature {
		feature[i] -= mean
	}

	sstotal := SumofSquares(feature)

	sswn := sstotal - ssbn
	dfbn := float64(len(groups)) - 1  
	dfwn := bigN - float64(len(groups))
	
	msb := ssbn / dfbn
	msw := sswn / dfwn

	return msb / msw
}

func GroupByClass(data []float64, class []int) [][]float64 {
	
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
func SquaresofSum (data []float64) float64 {
	sum := 0.0

	for _, val := range data {
		sum += val 
	}

	return sum * sum
}

// return a^2 + b^2 + ...
func SumofSquares (data []float64) float64 {
	sum := 0.0

	for _, val := range data {
		sum += val * val 
	} 

	return sum 
}



func Mean(lst []float64) float64 {

	a := 0.0

	for _, val := range lst {
		a += val
	}

	return a / float64(len(lst))
}