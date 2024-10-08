package handler

import (
	"fmt"
	"net/http"
	"os"

	"triple-s/config"
	"triple-s/utils"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	commandType, isValid, err := utils.ValidateURL(r.URL.Path)
	if !isValid {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusNotFound)
	}

	switch commandType {
	case config.HandlerBucketList:
		switch r.Method {
		case "GET": // List All Buckets Endpoint: "/"
			err := listOfBuckets(w)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	case config.HandlerBucket:
		bucketHandler(w, r)
	case config.HandlerObject:

	}
}

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
		fmt.Fprintf(os.Stderr, "Could not handle method: %v", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
