package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"triple-s/config"
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
		err := deleteObject(w, bucketName, objectKey)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func uploadObject(w http.ResponseWriter, r *http.Request, bucketName, objectKey string) error {
	if exists, err := utils.IsObjectExist(bucketName, objectKey); exists {
		http.Error(w, "Object already exists", http.StatusBadRequest)
		return err
	}

	bucketPath := filepath.Join(config.Dir, bucketName)
	// Create the destination file
	file := filepath.Join(bucketPath, objectKey)
	outFile, err := os.Create(file)
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

	// Get Content-Type and Content-Length
	contentType := r.Header.Get("Content-Type")
	contentLength := fmt.Sprint(r.ContentLength)

	// Adding new metadata to object.csv
	newBucketMetadata := []string{objectKey, contentLength, contentType, utils.GetCurrentTimeStamp()}
	err = utils.AddRowToCSV(filepath.Join(bucketPath, "objects.csv"), newBucketMetadata)
	if err != nil {
		return err
	}

	// Updating bucket.csv LastModifiedTime
	err = utils.UpdateField(filepath.Join(config.Dir, "buckets.csv"), bucketName, "LastModifiedTime", utils.GetCurrentTimeStamp())
	if err != nil {
		return err
	}

	fmt.Fprintln(w, "Object Created successfully")
	return nil
}

func retrieveObject(w http.ResponseWriter, bucketName, objectKey string) error {
	if exists, err := utils.IsObjectExist(bucketName, objectKey); !exists {
		http.Error(w, "Object does not exists", http.StatusBadRequest)
		return err
	}
	// Define the path to the object
	objectPath := filepath.Join(config.Dir, bucketName, objectKey)

	// Open the file
	file, err := os.Open(objectPath)
	if err != nil {
		http.Error(w, "Could not open object", http.StatusInternalServerError)
		return err
	}
	defer file.Close()

	// Set the Content-Type header
	metadata, err := utils.GetRow(filepath.Join(config.Dir, bucketName, "objects.csv"), "ObjectKey", objectKey)
	if err != nil {
		http.Error(w, "Could not open object", http.StatusInternalServerError)
		return err
	}
	contentLength := metadata[1]
	contentType := metadata[2]

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", contentLength)

	// Copy the file content to the response writer
	if _, err := io.Copy(w, file); err != nil {
		http.Error(w, "Could not send object", http.StatusInternalServerError)
		return err
	}
	return nil
}

func deleteObject(w http.ResponseWriter, bucketName, objectKey string) error {
	if exists, err := utils.IsObjectExist(bucketName, objectKey); !exists {
		http.Error(w, "Object does not exists", http.StatusBadRequest)
		return err
	}

	err := os.Remove(filepath.Join(config.Dir, bucketName, objectKey))
	if err != nil {
		http.Error(w, "Could not delete object", http.StatusInternalServerError)
		return err
	}
	// Updating object metadata
	err = utils.DeleteRow(filepath.Join(config.Dir, bucketName, "objects.csv"), objectKey)
	if err != nil {
		return err
	}

	// Updating Buckets metadata LastModifiedTime
	err = utils.UpdateField(filepath.Join(config.Dir, "buckets.csv"), bucketName, "LastModifiedTime", utils.GetCurrentTimeStamp())
	if err != nil {
		return err
	}

	fmt.Fprintln(w, "Object deleted successfully")
	return nil
}
