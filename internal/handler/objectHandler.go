package handler

import (
	"fmt"
	"net/http"
	"os"
)

func objectHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		err := uploadObject()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	case "GET":
		err := retrieveObject()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	case "DELETE":
		err := deleteObject()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func uploadObject() error {
	return nil
}

func retrieveObject() error {
	return nil
}

func deleteObject() error {
	return nil
}
