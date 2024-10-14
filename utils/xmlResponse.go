package utils

import (
	"encoding/xml"
	"net/http"
	"triple-s/internal/model"
)

func SendXmlListResponse(w http.ResponseWriter, response any) {
	w.Header().Set("Content-Type", "application/xml")
	// Encode the response to XML
	if err := xml.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode XML response", http.StatusInternalServerError)
	}
}

func SendXmlResponse(w http.ResponseWriter, response model.XMLResponse) {
	w.Header().Set("Content-Type", "application/xml")

	if response.Status != http.StatusOK {
		// Only write the header if it is not a 204 response
		if response.Status != http.StatusNoContent {
			w.WriteHeader(response.Status)
			if err := xml.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, "Failed to encode XML response", http.StatusInternalServerError)
			}
		} else {
			w.WriteHeader(response.Status)
		}
	} else {
		if err := xml.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode XML response", http.StatusInternalServerError)
		}
	}
}
