package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"triple-s/config"
	"triple-s/utils"
)

func bucketHandler(w http.ResponseWriter, r *http.Request) {
	bucketName := strings.TrimLeft(r.URL.Path, "/")
	switch r.Method {
	case http.MethodPut:
		err := createBucket(w, bucketName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	case http.MethodDelete:
		err := deleteBucket(w, bucketName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func createBucket(w http.ResponseWriter, bucketName string) error {
	bucketPath := filepath.Join(config.Dir, bucketName)
	isBucketExist, err := utils.SearchValueCSV(config.Dir+"/buckets.csv", "Name", bucketName)
	if err != nil {
		http.Error(w, "Error creating Bucket", http.StatusInternalServerError) // 500 Internal Server Error
		return err
	}

	if isBucketExist {
		http.Error(w, "Bucket already exists", http.StatusConflict) // 409 Conflict
		return config.ErrBucketExists
	}
	err = os.Mkdir(bucketPath, 0o755) // 0755/0700 is the permission mode
	if err != nil {
		http.Error(w, "Error creating Bucket", http.StatusInternalServerError) // 500 Internal Server Error
		return err
	}

	// Creating objects.csv for storing metadata
	newObjectsMetadataPath := filepath.Join(bucketPath, "/objects.csv")
	err = os.WriteFile(newObjectsMetadataPath, config.ObjectMetadataFields, 0o755)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating metadata: %v", err), http.StatusInternalServerError)
		return err
	}

	// Writing metadata to buckets.csv
	newBucketMetadata := []string{bucketName, utils.GetCurrentTimeStamp(), utils.GetCurrentTimeStamp(), "active"}
	err = utils.AddRowToCSV(filepath.Join(config.Dir, "/buckets.csv"), newBucketMetadata)
	if err != nil {
		return err
	}

	fmt.Fprint(w, "Bucket created successfully")
	return nil
}

func deleteBucket(w http.ResponseWriter, bucketName string) error {
	bucketPath := filepath.Join(config.Dir, bucketName)
	fmt.Println("Removing bucket", bucketPath)

	// Check if bucket exists
	metadataDir := filepath.Join(config.Dir, "/buckets.csv")
	isBucketExist, err := utils.SearchValueCSV(metadataDir, "Name", bucketName)
	if err != nil {
		http.Error(w, "Error creating Bucket", http.StatusInternalServerError) // 500 Internal Server Error
		return err
	}

	if !isBucketExist {
		http.Error(w, "Bucket does not exists", http.StatusConflict) // 409 Conflict
		return config.ErrBucketDoesNotExist
	}

	// Check if any object exist in this bucket
	col, err := utils.GetColumn(filepath.Join(bucketPath, "/objects.csv"), 0)
	if err != nil {
		http.Error(w, "Error deleting bucket", http.StatusInternalServerError)
		return err
	}
	if len(col) > 1 {
		http.Error(w, "Bucket is not empty", http.StatusBadRequest)
		return nil
	}

	// Removing bucket
	err = os.RemoveAll(bucketPath)
	if err != nil {
		http.Error(w, "Could not delete bucket", http.StatusInternalServerError)
		return err
	}

	// Updating metadata(buckets.csv)
	err = utils.DeleteRow(metadataDir, bucketName)
	if err != nil {
		fmt.Fprint(os.Stderr, "Error updating metadata, ERROR:")
		return err
	}

	fmt.Fprintf(w, "Bucket deleted successfully: %v", bucketName)
	return nil
}
