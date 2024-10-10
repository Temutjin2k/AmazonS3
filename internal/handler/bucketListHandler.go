package handler

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"triple-s/config"
	"triple-s/internal/model"
)

// List All Buckets Endpoint: "/"
func bucketListHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := listOfBuckets(w)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func listOfBuckets(w http.ResponseWriter) error {
	metadataDir := filepath.Join(config.Dir, "/buckets.csv")
	file, err := os.Open(metadataDir)
	if err != nil {
		http.Error(w, "Failed to open metadata file", http.StatusInternalServerError)
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll() // All records from buckets.csv
	if err != nil {
		http.Error(w, "Failed to read metadata file", http.StatusInternalServerError)
		return err
	}

	var buckets []model.Bucket
	for _, record := range records[1:] {
		if len(record) >= 4 {
			buckets = append(buckets, model.Bucket{
				Name:         record[0],
				CreationDate: record[1],
				LastModified: record[2],
				Status:       record[3],
			})
		} else {
			fmt.Println("in record less than 4 columns")
		}
	}
	response := model.BucketResponse{Buckets: buckets}
	w.Header().Set("Content-Type", "application/xml")

	// Encode the response to XML
	if err := xml.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode XML response", http.StatusInternalServerError)
	}
	return nil
}
