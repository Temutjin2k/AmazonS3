package handler

import (
	"fmt"
	"net/http"
	"os"

	"triple-s/config"
	"triple-s/utils"
)

func bucketHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT": // Create bucket Endpoint: "/{BucketName}"
		err := createBucket(w, r.URL.Path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	case "DELETE":
		err := deleteBucket(w, r.URL.Path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func createBucket(w http.ResponseWriter, urlPath string) error {
	bucketPath := config.Dir + urlPath
	bucketName := urlPath[1:]
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
	newObjectsMetadataPath := bucketPath + "/objects.csv"
	err = os.WriteFile(newObjectsMetadataPath, config.ObjectMetadataFields, 0o755)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating metadata: %v", err), http.StatusInternalServerError)
		return err
	}

	// Writing metadata to buckets.csv
	newBucketMetadata := []string{bucketName, utils.GetCurrentTimeStamp(), utils.GetCurrentTimeStamp(), "active"}
	err = utils.AddRowToCSV(config.Dir+"/buckets.csv", newBucketMetadata)
	if err != nil {
		return err
	}

	fmt.Fprint(w, "Bucket created successfully")
	return nil
}

func deleteBucket(w http.ResponseWriter, urlPath string) error {
	if urlPath == "/" {
		http.Error(w, "Error deleting bucket", http.StatusBadRequest)
		return config.ErrInvalidPath
	}
	bucketPath := config.Dir + urlPath
	fmt.Println("Removing bucket", bucketPath)

	// Check if path exists
	metadataDir := config.Dir + "/buckets.csv"
	isBucketExist, err := utils.SearchValueCSV(metadataDir, "Name", urlPath[1:])
	if err != nil {
		http.Error(w, "Error creating Bucket", http.StatusInternalServerError) // 500 Internal Server Error
		return err
	}

	if !isBucketExist {
		http.Error(w, "Bucket does not exists", http.StatusConflict) // 409 Conflict
		return config.ErrBucketDoesNotExist
	}

	info, err := os.Stat(bucketPath)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error checking bucket", http.StatusInternalServerError)
		return err
	}

	if !info.IsDir() {
		http.Error(w, fmt.Sprintf("%v is not Bucket", urlPath), http.StatusBadRequest)
		return err
	}

	// Check if any object exist in this bucket
	col, err := utils.GetColumn(bucketPath+"/objects.csv", 0)
	if err != nil {
		http.Error(w, "Error deleting bucket", http.StatusBadRequest)
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
	err = utils.DeleteRow(metadataDir, urlPath[1:])
	if err != nil {
		fmt.Fprint(os.Stderr, "Error updating metadata, ERROR:")
		return err
	}

	fmt.Fprintf(w, "Bucket deleted successfully: %v", urlPath[1:])
	return nil
}
