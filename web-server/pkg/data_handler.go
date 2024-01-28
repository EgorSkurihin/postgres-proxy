package web

import (
	"net/http"
	"os"
)

type DataHandler struct {
	dataFilePath string
}

func (h *DataHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile(h.dataFilePath)
	if err != nil {
		http.Error(w, "Unable to read data file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
