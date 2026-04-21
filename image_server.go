package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

func ImageServer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	camera := vars["camera"]
	day := vars["day"]
	hour := vars["hour"]
	image := vars["image"]

	filePath := filepath.Join("./camera-recordings", camera, day, "pic_001", hour, image+"[M][0@0][0].jpg")
	// check if file exists and is readable
	_, err := os.Stat(filePath)
	if err != nil {
		http.Error(w, "Bild nicht gefunden", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, filePath)
}
