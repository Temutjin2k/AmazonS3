package handler

import (
	"errors"
	"fmt"
	"net/http"
)

var ErrInvalidPath = errors.New("error deleting bucket")

func Handler(w http.ResponseWriter, r *http.Request) {
	// utils.ValidateURL(r.URL.Path)

	switch r.Method {
	case "PUT": // Create bucket Endpoint: "/{BucketName}"
		// TODO {BucketName} validation
		createBucket(w, r.URL.Path)
	case "GET": // List All Buckets Endpoint: "/"
		err := listOfBuckets(w)
		if err != nil {
			fmt.Println(err)
		}
	case "DELETE":
		err := deleteBucket(w, r.URL.Path)
		if err != nil {
			fmt.Println(err)
		}
	}
}
