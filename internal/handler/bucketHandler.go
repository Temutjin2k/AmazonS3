package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"triple-s/config"
	"triple-s/internal/model"
	"triple-s/utils"
)

func bucketHandler(w http.ResponseWriter, r *http.Request) {
	bucketName := strings.TrimLeft(r.URL.Path, "/")
	var response model.XMLResponse
	switch r.Method {
	case http.MethodPut:
		response = createBucket(bucketName)
	case http.MethodDelete:
		response = deleteBucket(bucketName)
	default:
		response.Status = http.StatusMethodNotAllowed
		response.Message = "Method not allowed"
	}

	response.Resource = r.URL.Path
	utils.SendXmlResponse(w, response)
}

func createBucket(bucketName string) model.XMLResponse {
	var response model.XMLResponse

	bucketPath := filepath.Join(config.Dir, bucketName)
	isBucketExist, err := utils.IsBucketExist(bucketName)
	if err != nil {
		fmt.Println(err)
		response.Status = http.StatusInternalServerError // 500
		response.Message = "Error creating Bucket"
		return response
	}

	if isBucketExist {
		response.Status = http.StatusBadRequest
		response.Message = "Bucket already exists"
		return response
	}

	err = os.Mkdir(bucketPath, 0o755) // 0755/0700 is the permission mode
	if err != nil {
		response.Status = http.StatusInternalServerError // 500
		response.Message = "Error creating Bucket"
		return response
	}

	// Creating objects.csv for storing metadata
	newObjectsMetadataPath := filepath.Join(bucketPath, "/objects.csv")
	err = os.WriteFile(newObjectsMetadataPath, config.ObjectMetadataFields, 0o755)
	if err != nil {
		response.Status = http.StatusInternalServerError // 500
		response.Message = "Error creating metadata"
		return response
	}

	// Writing metadata to buckets.csv
	newBucketMetadata := []string{bucketName, utils.GetCurrentTimeStamp(), utils.GetCurrentTimeStamp(), "active"}
	err = utils.AddRowToCSV(filepath.Join(config.Dir, "/buckets.csv"), newBucketMetadata)
	if err != nil {
		response.Status = http.StatusInternalServerError // 500
		response.Message = "Error creating metadata"
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Bucket Created successfully"
	return response
}

func deleteBucket(bucketName string) model.XMLResponse {
	response := model.XMLResponse{}

	bucketPath := filepath.Join(config.Dir, bucketName)

	// Check if bucket exists
	metadataDir := filepath.Join(config.Dir, "/buckets.csv")
	isBucketExist, err := utils.IsBucketExist(bucketName)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error creating Bucket"
		return response
	}

	if !isBucketExist {
		response.Status = http.StatusBadRequest
		response.Message = "Bucket does not exists"
		return response
	}

	// Check if any object exist in this bucket
	col, err := utils.GetColumn(filepath.Join(bucketPath, "/objects.csv"), 0)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error deleting bucket"
		return response
	}
	if len(col) > 1 {
		response.Status = http.StatusBadRequest
		response.Message = "Bucket is not empty"
		return response
	}

	// Removing bucket
	err = os.RemoveAll(bucketPath)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error deleting bucket"
		return response
	}

	// Updating metadata(buckets.csv)
	err = utils.DeleteRow(metadataDir, bucketName)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Error updating metadata"
		return response

	}

	response.Status = http.StatusOK
	response.Message = "Bucket deleted successfully"

	return response
}
