package handler

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"

	"triple-s/model"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT": // Create bucket Endpoint: "/{BucketName}"
		// TODO {BucketName} validation
		createBucket(w, r.URL.Path)
	case "GET": // List All Buckets Endpoint: "/"
		err := listOfBuckets(w)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func createBucket(w http.ResponseWriter, urlPath string) {
	folderName := "./data" + urlPath
	err := os.Mkdir(folderName, 0o700) // 0755/0700 is the permission mode
	if err != nil {
		if os.IsExist(err) {
			http.Error(w, "Folder already exists", http.StatusConflict) // 409 Conflict
		} else {
			http.Error(w, fmt.Sprintf("Error creating folder: %v", err), http.StatusInternalServerError) // 500 Internal Server Error
		}
	} else {
		fmt.Fprintf(w, "Folder created successfully: %v", folderName)
	}
}

func listOfBuckets(w http.ResponseWriter) error {
	file, err := os.Open("data/buckets.csv")
	if err != nil {
		http.Error(w, "Failed to open CSV file", http.StatusInternalServerError)
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll() // All records from buckets.csv
	if err != nil {
		http.Error(w, "Failed to read CSV file", http.StatusInternalServerError)
		return err
	}

	var buckets []model.Bucket
	for _, record := range records[1:] {
		if len(record) >= 3 {
			buckets = append(buckets, model.Bucket{
				CreationDate: record[0],
				Name:         record[1],
				LastModified: record[2],
			})
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
