package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

func ConvertStart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	camera := vars["camera"]
	day := vars["day"]
	video := vars["video"]
	hour := vars["hour"]
	clip := vars["clip"]

	clipID := strings.ReplaceAll(clip, "[M][0@0][0]", "")

	w.Header().Set("Content-Type", "text/html")
	responseHTML := fmt.Sprintf(`
		<div id="vid-%s" style="background-color: #fff3cd; color: #856404; padding: 10px; margin-bottom: 10px; border: 1px solid #ffeeba; border-radius: 5px;"
			 hx-get="/convert/%s/%s/%s/%s/%s.mp4" hx-trigger="load" hx-swap="outerHTML">
			⏳ Konvertiere Video <strong>%s</strong>... Bitte warten, dies kann einen Moment dauern.
		</div>`, clipID, camera, day, video, hour, clip, clip)
	fmt.Fprint(w, responseHTML)
}

func ClipConverter(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	camera := vars["camera"]
	day := vars["day"]
	video := vars["video"]
	hour := vars["hour"]
	clip := vars["clip"]

	clipID := strings.ReplaceAll(clip, "[M][0@0][0]", "")

	baseDirRecordings := "./camera-recordings"
	inputPath := filepath.Join(baseDirRecordings, camera, day, video, hour, clip+".dav")
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		inputPath = filepath.Join(baseDirRecordings, camera, day, video, hour, clip+".dav_")
	}

	baseDirConverted := "./converted"
	outputDir := filepath.Join(baseDirConverted, camera, day, hour)
	outputPath := filepath.Join(outputDir, clip+".mp4")

	fmt.Printf("converting %s\n", outputPath)

	os.MkdirAll(outputDir, os.ModePerm)

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		cmd := exec.Command("ffmpeg", "-y", "-i", inputPath, "-c:v", "copy", "-tag:v", "hvc1", outputPath)
		fmt.Println("converting")
		err := cmd.Run()
		if err != nil {
			http.Error(w, fmt.Sprintf("Fehler bei der Konvertierung: %v", err), http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html")
	responseHTML := fmt.Sprintf(`
		<div style="background-color: #d4edda; color: #155724; padding: 10px; margin-bottom: 10px; border: 1px solid #c3e6cb; border-radius: 5px;">
			✅ Clip erfolgreich konvertiert!
		</div>
		<video id="vid-%s" width="100%%" controls autoplay>
			<source src="/converted/%s/%s/%s/%s.mp4" type="video/mp4">
		</video>`, clipID, camera, day, hour, clip)
	fmt.Fprint(w, responseHTML)
}
