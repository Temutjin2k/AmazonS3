package handler

import (
	"fmt"
	"net/http"
	"os"

	"triple-s/config"
	"triple-s/internal/model"
	"triple-s/utils"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" && r.Method != "GET" && r.Method != "DELETE" {
		response := model.XMLResponse{
			Status:   http.StatusMethodNotAllowed,
			Message:  "Method not allowed",
			Resource: r.URL.Path,
		}
		utils.SendXmlResponse(w, response)
		return
	}
	commandType, isValid, err := utils.ValidateURL(r.URL.Path)
	if !isValid {
		fmt.Fprintln(os.Stderr, "Error Validating URL:", err)
		response := model.XMLResponse{
			Status:   http.StatusNotFound,
			Message:  "Not Found",
			Resource: r.URL.Path,
		}
		utils.SendXmlResponse(w, response)
		return
	}

	switch commandType {
	case config.HandlerBucketList:
		bucketListHandler(w, r)
	case config.HandlerBucket:
		bucketHandler(w, r)
	case config.HandlerObject:
		objectHandler(w, r)
	}
}
