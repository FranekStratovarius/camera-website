package main

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

func ClipServer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	camera := vars["camera"]
	day := vars["day"]
	video := vars["video"]
	hour := vars["hour"]
	clip := vars["clip"]

	filePath := filepath.Join("./camera-recordings", camera, day, video, hour, clip+".mp4")
	http.ServeFile(w, r, filePath)
}
