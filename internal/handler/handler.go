package handler

import (
	"fmt"
	"net/http"

	"triple-s/config"
	"triple-s/utils"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" && r.Method != "GET" && r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	commandType, isValid, err := utils.ValidateURL(r.URL.Path)
	if !isValid {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusNotFound)
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
