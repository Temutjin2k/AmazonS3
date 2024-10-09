package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"triple-s/utils"
)

func objectHandler(w http.ResponseWriter, r *http.Request) {
	// Get the bucket name and object key from the URL
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 2 {
		http.NotFound(w, r)
		return
	}
	bucketName := parts[0]
	objectKey := parts[1]
	if exists, err := utils.IsBucketExist(bucketName); !exists {
		fmt.Fprintln(os.Stderr, err)
		http.Error(w, "Bucket does not exists", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPut:
		err := uploadObject(w, r, bucketName, objectKey)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	case http.MethodGet:
		err := retrieveObject(w, bucketName, objectKey)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	case http.MethodDelete:
		err := deleteObject()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func uploadObject(w http.ResponseWriter, r *http.Request, bucketName, objectKey string) error {
	if exists, err := utils.IsObjectExist(bucketName, objectKey); exists {
		fmt.Fprintln(os.Stderr, err)
		http.Error(w, "Object already exists", http.StatusBadRequest)
		return err
	}

	bucketPath := filepath.Join("./data", bucketName)
	// Create the destination file
	filePath := filepath.Join(bucketPath, objectKey)
	outFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Could not create object", http.StatusInternalServerError)
		return err
	}
	defer outFile.Close()

	// Copy the uploaded file data to the destination file
	_, err = io.Copy(outFile, r.Body)
	if err != nil {
		http.Error(w, "Could not save object", http.StatusInternalServerError)
		return err
	}

	fmt.Fprintln(w, "Object Created successfully")
	return nil
}

func retrieveObject(w http.ResponseWriter, bucketName, objectKey string) error {
	if exists, err := utils.IsObjectExist(bucketName, objectKey); !exists {
		fmt.Fprintln(os.Stderr, err)
		http.Error(w, "Object does not already exists", http.StatusBadRequest)
		return err
	}
	return nil
}

func deleteObject() error {
	return nil
}
