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
	case "DELETE":
		if r.URL.Path == "" {
			http.Error(w, "Error deleting bucket", http.StatusInternalServerError)
		}
		path := "./data" + r.URL.Path
		fmt.Println("Removing directory", path)

		err := os.RemoveAll(path)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintf(w, "Bucket deleted successfully: %v", path)
	}
}

func createBucket(w http.ResponseWriter, urlPath string) {
	bucketName := "./data" + urlPath
	err := os.Mkdir(bucketName, 0o755) // 0755/0700 is the permission mode
	if err != nil {
		if os.IsExist(err) {
			http.Error(w, "Bucket already exists", http.StatusConflict) // 409 Conflict
		} else {
			http.Error(w, fmt.Sprintf("Error creating Bucket: %v", err), http.StatusInternalServerError) // 500 Internal Server Error
		}
	} else {
		fmt.Fprintf(w, "Bucket created successfully: %v", bucketName)
	}
}

func listOfBuckets(w http.ResponseWriter) error {
	file, err := os.Open("data/buckets.csv")
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
