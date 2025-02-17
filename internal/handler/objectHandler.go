package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"triple-s/config"
	"triple-s/internal/model"
	"triple-s/utils"
)

func objectHandler(w http.ResponseWriter, r *http.Request) {
	// var response model.XMLResponse
	var response model.XMLResponse
	var err error
	response.Status = http.StatusInternalServerError
	response.Message = "Server Error"
	response.Resource = r.URL.Path

	// Get the bucket name and object key from the URL
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 2 {
		response := model.XMLResponse{
			Status:  http.StatusNotFound,
			Message: "Not Found",
		}
		utils.SendXmlResponse(w, response)
		return
	}
	bucketName := parts[0]
	objectKey := parts[1]
	if exists := utils.IsBucketExist(bucketName); !exists {
		response := model.XMLResponse{
			Status:  http.StatusBadRequest,
			Message: "Bucket does not exists",
		}
		utils.SendXmlResponse(w, response)
		return
	}

	switch r.Method {
	case http.MethodPut:
		response, err = uploadObject(r, bucketName, objectKey)
	case http.MethodGet:
		err := retrieveObject(w, bucketName, objectKey)
		if err == config.ErrObjectDoesNotExist {
			response.Status = http.StatusBadRequest
			response.Message = "Object does not exist"
			utils.SendXmlResponse(w, response)
		} else if err != nil {
			fmt.Fprintln(os.Stderr, err)
			response.Status = http.StatusInternalServerError
			response.Message = "Error retriving object"
			utils.SendXmlResponse(w, response)
		}

		return
	case http.MethodDelete:
		response, err = deleteObject(bucketName, objectKey)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	utils.SendXmlResponse(w, response)
}

func uploadObject(r *http.Request, bucketName, objectKey string) (model.XMLResponse, error) {
	var response model.XMLResponse

	bucketPath := filepath.Join(config.Dir, bucketName)

	exists := utils.IsObjectExist(bucketName, objectKey)

	// Updating metadata
	if exists {
		err := utils.DeleteRow(filepath.Join(bucketPath, "objects.csv"), objectKey)
		if err != nil {
			response.Status = http.StatusInternalServerError
			response.Message = "Error Updating metadata"
			return response, err
		}
	}

	// The destination file
	file := filepath.Join(bucketPath, objectKey)
	outFile, err := os.Create(file)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error creating object"
		return response, err
	}
	defer outFile.Close()

	// Copy the uploaded file data to the destination file
	_, err = io.Copy(outFile, r.Body)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error creating object"
		return response, err
	}

	// Get Content-Type and Content-Length
	contentType := r.Header.Get("Content-Type")
	contentLength := fmt.Sprint(r.ContentLength)

	// Adding new metadata to object.csv
	newBucketMetadata := []string{objectKey, contentLength, contentType, utils.GetCurrentTimeStamp()}
	err = utils.AddRowToCSV(filepath.Join(bucketPath, "objects.csv"), newBucketMetadata)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error creating object metadata"
		return response, err
	}

	// Updating bucket.csv LastModifiedTime
	err = utils.UpdateField(filepath.Join(config.Dir, "buckets.csv"), bucketName, "LastModifiedTime", utils.GetCurrentTimeStamp())
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error updating bucket's metadata"
		return response, err
	}

	response.Status = http.StatusOK
	response.Message = "Object uploaded successfully"
	return response, nil
}

func retrieveObject(w http.ResponseWriter, bucketName, objectKey string) error {
	if exists := utils.IsObjectExist(bucketName, objectKey); !exists {
		err := config.ErrObjectDoesNotExist
		return err
	}
	// Define the path to the object
	objectPath := filepath.Join(config.Dir, bucketName, objectKey)

	// Open the file
	file, err := os.Open(objectPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Set the Content-Type header
	metadata, err := utils.GetRow(filepath.Join(config.Dir, bucketName, "objects.csv"), "ObjectKey", objectKey)
	if err != nil {
		return err
	}
	contentLength := metadata[1]
	contentType := metadata[2]

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", contentLength)

	// Copy the file content to the response writer
	if _, err := io.Copy(w, file); err != nil {
		return err
	}

	return nil
}

func deleteObject(bucketName, objectKey string) (model.XMLResponse, error) {
	var response model.XMLResponse

	if exists := utils.IsObjectExist(bucketName, objectKey); !exists {
		response.Status = http.StatusBadRequest
		response.Message = "Object does not exists"
		return response, config.ErrObjectDoesNotExist
	}

	// Deleting object
	err := os.Remove(filepath.Join(config.Dir, bucketName, objectKey))
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error deleting object"
		return response, err
	}

	// Updating object metadata
	err = utils.DeleteRow(filepath.Join(config.Dir, bucketName, "objects.csv"), objectKey)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error Updating metadata"
		return response, err
	}

	// Updating Buckets metadata LastModifiedTime
	err = utils.UpdateField(filepath.Join(config.Dir, "buckets.csv"), bucketName, "LastModifiedTime", utils.GetCurrentTimeStamp())
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error Updating metadata"
		return response, err
	}

	response.Status = http.StatusNoContent
	response.Message = "Object deleted successfully"
	return response, nil
}
