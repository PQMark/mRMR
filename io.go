package mRMR

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"sort"
)

// return Data, Features, and Class
func ReadCSV(filepath string, irrelevantCols, irrelevantRows []int, featureIndex, groupIndex int, colFeatures bool) ([][]float64, []string, []int) {
	featureIndex -= 1
	groupIndex -= 1

	irrelevantCols = convertToZeroBased(irrelevantCols)
	irrelevantRows = convertToZeroBased(irrelevantRows)

	file, err := os.Open(filepath)
	if err != nil {
		panic(fmt.Sprintf("unable to open file %s: %v", filepath, err))
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		panic(fmt.Sprintf("error reading CSV data: %v", err))
	}

	// Remove irrelevant rows
	if len(irrelevantRows) > 0 {
		records = removeRows(records, irrelevantRows)
	}

	// Remove irrelevant columns
	if len(irrelevantCols) > 0 {
		records = removeColumns(records, irrelevantCols)
	}

	// Transpose if colFeatures is false
	if !colFeatures {
		records = transpose(records)
	}

	var groups []int
	groupMap := make(map[string]int)
	currentGroup := 0

	if groupIndex < 0 {
		panic("groupIndex cannot be negative")
	}
	for i, row := range records {
		if groupIndex >= len(row) {
			panic(fmt.Sprintf("groupIndex %d out of range in row %d", groupIndex, i))
		}
		groupStr := row[groupIndex]
		if groupStr == "NA" {
			panic(fmt.Sprintf("NA encountered at groupIndex %d in row %d", groupIndex, i))
		}

		if _, exists := groupMap[groupStr]; !exists {
			groupMap[groupStr] = currentGroup
			currentGroup++
		}
		groups = append(groups, groupMap[groupStr])
	}

	// assuming the first element is the name
	groups = groups[1:]
	records = removeColumns(records, []int{groupIndex})

	var features []string
	if featureIndex != -1 {
		if featureIndex < 0 || featureIndex >= len(records) {
			panic(fmt.Sprintf("featureIndex %d out of range", featureIndex))
		}
		features = records[featureIndex]
		
		records = removeRows(records, []int{featureIndex})
	}

	data := make([][]float64, len(records))
	for i, row := range records {
		data[i] = make([]float64, len(row))
		for j, field := range row {
			if field == "NA" {
				panic(fmt.Sprintf("NA encountered at row %d, column %d", i, j))
			}
			num, err := strconv.ParseFloat(field, 64)
			if err != nil {
				panic(fmt.Sprintf("invalid float at row %d, column %d: %v", i, j, err))
			}
			data[i][j] = num
		}
	}

	return data, features, groups
}

func removeRows(records [][]string, irrelevantRows []int) [][]string {
	sort.Sort(sort.Reverse(sort.IntSlice(irrelevantRows)))
	for _, idx := range irrelevantRows {
		if idx < 0 || idx >= len(records) {
			panic(fmt.Sprintf("irrelevant row index %d out of range", idx))
		}
		records = append(records[:idx], records[idx+1:]...)
	}
	return records
}

func removeColumns(records [][]string, irrelevantCols []int) [][]string {
	sort.Sort(sort.Reverse(sort.IntSlice(irrelevantCols)))
	for _, idx := range irrelevantCols {
		for i, row := range records {
			if idx < 0 || idx >= len(row) {
				panic(fmt.Sprintf("irrelevant column index %d out of range in row %d", idx, i))
			}
			records[i] = append(row[:idx], row[idx+1:]...)
		}
	}
	return records
}

func transpose(matrix [][]string) [][]string {
	if len(matrix) == 0 {
		return matrix
	}
	transposed := make([][]string, len(matrix[0]))
	for i := range transposed {
		transposed[i] = make([]string, len(matrix))
		for j := range matrix {
			transposed[i][j] = matrix[j][i]
		}
	}
	return transposed
}

func convertToZeroBased(indices []int) []int {
	zeroBased := make([]int, len(indices))
	for i, idx := range indices {
		if idx < 1 {
			panic(fmt.Sprintf("invalid index %d: must be >=1", idx))
		}
		zeroBased[i] = idx - 1
	}
	return zeroBased
}