package utils

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"
)

// Search value in CSV file, return rows index and boolean isExist. Returns (-1, false) if target string was not found. This func skips first row
func SearchValueCSV(filepath string, colName string, target string) (bool, error) {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		return false, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		return false, err
	}

	if len(records) == 0 {
		return false, errors.New("not enough rows in file")
	}

	targetCol := -1
	for i, v := range records[0] {
		if v == colName {
			targetCol = i
			break
		}
	}

	if targetCol == -1 {
		return false, errors.New("could not find given column")
	}
	for _, row := range records[1:] {
		if row[targetCol] == target {
			return true, nil
		}
	}

	return false, nil
}

// func SearchValueCSV(filepath string, col int, target string) bool {
// 	file, err := os.Open(filepath)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Error: %v", err)
// 		return false
// 	}
// 	defer file.Close()
// 	reader := csv.NewReader(file)
// 	records, err := reader.ReadAll()
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Error: %v", err)
// 		return false
// 	}

// 	for _, row := range records[1:] {
// 		if row[col] == target {
// 			return true
// 		}
// 	}

// 	return false
// }

func AddRowToCSV(filename string, record []string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(record); err != nil {
		return err
	}

	return nil
}

// Delete row by compating target value with first column
func DeleteRow(filename string, target string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	var updatedRecords [][]string

	// Filter out the rows that have the target in the first column
	for _, record := range records {
		if len(record) > 0 && record[0] != target {
			updatedRecords = append(updatedRecords, record)
		}
	}

	// For Test
	if len(updatedRecords) == len(records) {
		return errors.New("no matching row found to delete")
	}

	file, err = os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the updated records back to the file
	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.WriteAll(updatedRecords); err != nil {
		return err
	}
	return nil
}

func GetColumn(filepath string, col int) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var columnValues []string

	for _, record := range records {
		if col < len(record) && strings.ReplaceAll(record[col], " ", "") != "" {
			columnValues = append(columnValues, record[col])
		} else {
			return nil, fmt.Errorf("column index %d out of range", col)
		}
	}

	return columnValues, nil
}
