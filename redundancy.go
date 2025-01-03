package mRMR

import "math"

// RedundancyUpdate calculates the redundancy between each unselected feature with last selected feature and updates the redundancy map.
func RedundancyUpdate(data [][]float64, featureToConsider []int, target int, redundancyMap map[[2]int]float64, redundancyFunc func([]float64, []float64) float64) map[[2]int]float64 {

	data2 := getCol(data, target)

	for _, idx := range featureToConsider {
		data1 := getCol(data, idx)
		mi := redundancyFunc(data1, data2)
		redundancyMap[[2]int{target, idx}] = mi
	}

	return redundancyMap
}

// PearsonCorrelation returns the absolute value of pearson correlation coefficient
func PearsonCorrelation(data1, data2 []float64) float64 {
	if len(data1) != len(data2) {
		panic("feature slices must have the same length")
	}
	
	mean1 := mean(data1)
	mean2 := mean(data2)

	for i := range data1 {
		data1[i] -= mean1
		data2[i] -= mean2
	}

	sd1 := sumOfSquares(data1)
	sd1 = math.Sqrt(sd1)

	sd2 := sumOfSquares(data2)
	sd2 = math.Sqrt(sd2)

	cov := 0.0
	for i := range data1 {
		cov += data1[i] * data2[i]
	}

	return math.Abs(cov / (sd1 * sd2))
}