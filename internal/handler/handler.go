package handler

import (
	"encoding/csv"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"os"

	"triple-s/internal/model"
)

var ErrInvalidPath = errors.New("error deleting bucket")

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

		err := deleteBucket(w, r.URL.Path)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func createBucket(w http.ResponseWriter, urlPath string) {
	bucketName := "./data" + urlPath
	err := os.Mkdir(bucketName, 0o755) // 0755/0700 is the permission mode
	if err != nil {
		if os.IsExist(err) {
			http.Error(w, "Bucket already exists", http.StatusConflict) // 409 Conflict
			return
		} else {
			http.Error(w, fmt.Sprintf("Error creating Bucket: %v", err), http.StatusInternalServerError) // 500 Internal Server Error
			return
		}
	} else {
		fmt.Fprint(w, "Bucket created successfully")
	}

	// Creating objects.csv for storing metadata
	fmt.Println(bucketName)
	err = os.WriteFile(bucketName+"/objects.csv", []byte("ObjectKey,Size,ContentType,LastModified"), 0o755)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating metadata: %v", err), http.StatusInternalServerError)
		fmt.Println(err)
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

func deleteBucket(w http.ResponseWriter, urlPath string) error {
	if urlPath == "/" {
		http.Error(w, "Error deleting bucket", http.StatusBadRequest)
		return ErrInvalidPath
	}
	path := "./data" + urlPath
	fmt.Println("Removing directory", path)
	// Check if path exists
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		http.Error(w, "Path does not exist", http.StatusNotFound)
		return err
	} else if err != nil {
		fmt.Println(err)
		http.Error(w, "Error checking path", http.StatusInternalServerError)
		return err
	}

	if !info.IsDir() {
		http.Error(w, "Path is not a directory", http.StatusBadRequest)
		return err
	}

	err = os.RemoveAll(path)
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
		return err
	}
	fmt.Fprintf(w, "Bucket deleted successfully: %v", path)
	return nil
}
