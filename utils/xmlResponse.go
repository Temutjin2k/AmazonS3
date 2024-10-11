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
		w.WriteHeader(response.Status)
	}
	// Encode the response to XML
	if err := xml.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode XML response", http.StatusInternalServerError)
	}
}
