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
	Data				DatamRMR
	Discretization		bool 
	BinSize				int
	Method 				string
	Calculation 		string 
	MaxFeatures			int	
	RedundancyMethod 	string
	Threshold			float64
	QLevel				int
	RelevanceFunc		func ([]float64, []int) float64
	RedundancyFunc  	func ([]float64, []float64) float64
}

type DatamRMR struct{
	X 	[][]float64
	Class []int
}

func (paras *ParasmRMR) MRMR() ([]int, []float64, map[[2]int]float64){
	
	paras.defaults()
	paras.setups()

	if paras.Discretization && paras.Method != "nmi-nmi"{
		paras.Data.X, _ = Discretization(paras.Data.X, paras.BinSize)
	}

	if paras.Method == "nmi-nmi" {
		_, paras.Data.X = Discretization(paras.Data.X, paras.QLevel)
	}

	relevanceAll := Relevance(paras.Data.X, paras.Data.Class, paras.RelevanceFunc)

	// Filter out features with zero relevance
	featuresToConsider := make([]int, 0, len(paras.Data.X[0]))
	for i, val := range relevanceAll {
		if val > 0 {
			featuresToConsider = append(featuresToConsider, i)
		}
	}
	
	if paras.Method == "fs-pearson" {    //fs
		relevanceAll = MinMaxNormalization(relevanceAll)
	}

	if paras.Method == "nmi-nmi" {
		n := UniqueClass(paras.Data.Class)
		relevanceAll = Scaling(relevanceAll, math.Log2(float64(n)))
	}

	selectedFeatures := make([]int, 0, paras.MaxFeatures)
	redundancyMap := make(map[[2]int]float64)

	for c := 0; c < paras.MaxFeatures; c++ {

		relevance := selectByIndex(relevanceAll, featuresToConsider)
		redundancy := make([]float64, len(featuresToConsider))

		if c != 0 {
			// calculate redundancy
			lastSelectedF := selectedFeatures[len(selectedFeatures) - 1]

			// update map 
			redundancyMap = RedundancyUpdate(paras.Data.X, featuresToConsider, lastSelectedF, redundancyMap, paras.RedundancyFunc)

			for i, f := range featuresToConsider {
				s := 0.0
				for _, sel := range selectedFeatures {
					val := redundancyMap[[2]int{sel, f}]

					if paras.RedundancyMethod == "avg"{
						s += val
					} else if paras.RedundancyMethod == "max" {
						if val > s {
							s = val
						}
					}
				}

				if paras.RedundancyMethod == "avg" {
					divisor := float64(len(selectedFeatures))
					if paras.Method == "nmi-nmi" {
						divisor = math.Log2(float64(paras.QLevel))
					} 

					redundancy[i] = s / divisor
				}

				if paras.RedundancyMethod == "max" {
					redundancy[i] = s 
				}
				
			}
			
		}

		fmt.Println("relevance:", relevance)
		fmt.Println("Redundancy:", redundancy)

		score := PairwiseOperation(relevance, redundancy, paras.Calculation)
		fmt.Println(score, "\n")

		// Early stopping
		if (paras.Calculation == "diff" && CheckIfAllNegative(score)) ||
			(paras.Calculation == "quo" && CheckIfAllSmallerOne(score)) {
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
	if paras.BinSize == 0 {
		paras.BinSize = int(math.Sqrt(float64(len(paras.Data.X))))
	}

	if paras.Calculation == "" {
		paras.Calculation = "diff"
	}

	if paras.Method == "" {
		paras.Method = "mi-mi"
	}
	
	if paras.RedundancyMethod == "" {
		paras.RedundancyMethod = "avg"
	}

	if paras.MaxFeatures == 0 {
		paras.MaxFeatures = len(paras.Data.X[0])
	}

	if paras.Threshold == 0{
		paras.Threshold = 0.01
	}

	if paras.MaxFeatures > len(paras.Data.X[0]) {
		log.Printf("Warning: maxFeatures (%d) exceeds number of features (%d). Adjusting.",
			paras.MaxFeatures, len(paras.Data.X[0]))
		paras.MaxFeatures = len(paras.Data.X[0])
	}

}

func (paras *ParasmRMR) setups() {
	
	switch strings.ToLower(paras.Method) {
	case "mi-mi":
		paras.RelevanceFunc = MutualInfo
		paras.RedundancyFunc = MutualInfo
	case "fs-pearson":
		paras.RelevanceFunc = FStatistic
		paras.RedundancyFunc = PearsonCorrelation
	case "nmi-nmi":
		paras.RelevanceFunc = MutualInfo
		paras.RedundancyFunc = MutualInfo
		paras.QLevel = QuantizationLevel(paras.Data.X, paras.Threshold)
	default:
		panic("Invalid method. Choose from 'mi-mi', 'fs-pearson', 'nmi-nmi'")
	}
}