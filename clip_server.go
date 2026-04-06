package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

func ClipServer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("serve video file")
	// fmt.Printf("w: %+v\n", w)
	// fmt.Printf("r: %+v\n", r)

	vars := mux.Vars(r)

	camera := vars["camera"]
	day := vars["day"]
	video := vars["video"]
	hour := vars["hour"]
	clip := vars["clip"]

	baseDir := "./camera-recordings"

	filePath := filepath.Join(baseDir, camera, day, video, hour, clip)
	fmt.Printf("VideoPath: %s\n", filePath)
	// http.ServeFile(w, r, filePath)
	http.ServeFile(w, r, "./converted/out.mp4")
}
