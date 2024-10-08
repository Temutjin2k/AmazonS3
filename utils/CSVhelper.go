package utils

import (
	"encoding/csv"
	"fmt"
	"os"
)

// Search value in CSV file, return rows index and boolean isExist. Returns (-1, false) if target string was not found. This func skips first row
func SearchValueCSV(filepath string, col int, target string) bool {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		return false
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		return false
	}

	for _, row := range records[1:] {
		if row[col] == target {
			return true
		}
	}

	return false
}

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
