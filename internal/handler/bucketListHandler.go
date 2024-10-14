package handler

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"triple-s/config"
	"triple-s/internal/model"
	"triple-s/utils"
)

// List All Buckets Endpoint: "/"
func bucketListHandler(w http.ResponseWriter, r *http.Request) {
	var response model.XMLResponse
	response.Resource = r.URL.Path

	switch r.Method {
	case http.MethodGet:
		response, err := listOfBuckets(w)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			utils.SendXmlResponse(w, response)
		}
	default:
		response.Status = http.StatusMethodNotAllowed
		response.Message = "Method not allowed"
		utils.SendXmlResponse(w, response)
	}
}

func listOfBuckets(w http.ResponseWriter) (model.XMLResponse, error) {
	var response model.XMLResponse

	metadataDir := filepath.Join(config.Dir, "/buckets.csv")
	file, err := os.Open(metadataDir)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Failed to open metadata file"
		return response, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll() // All records from buckets.csv
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = "Failed to open metadata file"
		return response, err
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
			fmt.Fprintln(os.Stderr, "ListOFBuckets: in record less than 4 columns")
		}
	}
	// Responsing buckets
	bucketResponse := model.BucketResponse{Buckets: buckets}
	utils.SendXmlListResponse(w, bucketResponse)

	return model.XMLResponse{}, nil
}
