package mRMR

import (
	"math"
	"strings"
	"fmt"
	"log"
)

type Numeric interface{
	int | int8 | int16 | int32 | int64 | float32 | float64
}

type ParasmRMR struct{
	Data			DatamRMR
	Discretization	bool 
	binSize			int
	Relevance 		string 
	Redundancy 		string
	Calculation 	string 
	maxFeatures		int	
	RelevanceFunc	func ([]float64, []int) float64
	RedundancyFunc  func ([]float64, []float64) float64
}

type DatamRMR struct{
	X 	[][]float64
	class []int
}

func (paras *ParasmRMR) mRMR() ([]int, []float64, map[[2]int]float64){
	
	paras.defaults()

	if paras.Discretization {
		paras.Data.X = Discretization(paras.Data.X, paras.binSize)
	}

	relevanceAll := Relevance(paras.Data.X, paras.Data.class, paras.RelevanceFunc)

	featuresToConsider := make([]int, 0, len(paras.Data.X[0]))	
	selectedFeatures := make([]int, 0, len(paras.Data.X[0]))
	redundancyMap := make(map[[2]int]float64)

	// Discard if relevance MI = 0
	for i, val := range relevanceAll {
		if val > 0.0 {
			featuresToConsider = append(featuresToConsider, i)
		}
	}
	
	for c := 0; c < paras.maxFeatures; c++ {

		relevance := selectByIndex(relevanceAll, featuresToConsider)
		redundancy := make([]float64, len(featuresToConsider))

		if c != 0 {
			// calculate redundancy
			lastSelectedF := selectedFeatures[len(selectedFeatures) - 1]

			// update map 
			redundancyMap = RedundancyUpdate(paras.Data.X, featuresToConsider, lastSelectedF, redundancyMap, paras.RedundancyFunc)

			for i, val1 := range featuresToConsider {
				sum := 0.0

				for _, val2 := range selectedFeatures {
					key := [2]int{val2, val1}
					sum += redundancyMap[key]
				}

				redundancy[i] = sum / float64(len(selectedFeatures))
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

func (paras *ParasmRMR) defaults() {
	if paras.binSize == 0 {
		paras.binSize = int(math.Sqrt(float64(len(paras.Data.X))))
	}

	if paras.Calculation == "" {
		paras.Calculation = "diff"
	}

	if paras.Relevance == "" {
		paras.Relevance = "mutual"
	}

	if paras.Redundancy == "" {
		paras.Redundancy = "mutual"
	}
	
	if paras.maxFeatures == 0 {
		paras.maxFeatures = 20
	}

	if paras.Relevance == strings.ToLower("MI") {
		paras.RedundancyFunc = MutualInfo
	}

	// check the methods

	if paras.maxFeatures > len(paras.Data.X[0]) {
		log.Printf("Warning: maxFeatures value exceeds the number of features. Adjusting maxFeatures to the number of features.")
		paras.maxFeatures = len(paras.Data.X[0])
	}
}