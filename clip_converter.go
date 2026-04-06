package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gorilla/mux"
)

type ConvertingPageData struct {
	CameraName string
	DayName    string
	VideoName  string
	HourName   string
	ClipName   string
}

func ConvertStart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	camera := vars["camera"]
	day := vars["day"]
	video := vars["video"]
	hour := vars["hour"]
	clip := vars["clip"]

	// clipID := strings.ReplaceAll(clip, "[M][0@0][0]", "")

	w.Header().Set("Content-Type", "text/html")

	data := ConvertingPageData{
		CameraName: camera,
		DayName:    day,
		VideoName:  video,
		HourName:   hour,
		ClipName:   clip,
	}

	err := tmpl.ExecuteTemplate(w, "converting.html", data)
	CheckTemplateError(w, err)
}

func ClipConverter(w http.ResponseWriter, r *http.Request) {
	fmt.Println("converting")
	vars := mux.Vars(r)

	camera := vars["camera"]
	day := vars["day"]
	video := vars["video"]
	hour := vars["hour"]
	clip := vars["clip"]

	// clipID := strings.ReplaceAll(clip, "[M][0@0][0]", "")

	clipDirectory := filepath.Join("./camera-recordings", camera, day, video, hour)

	inputPath := filepath.Join(clipDirectory, clip+".dav")
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		inputPath = filepath.Join(clipDirectory, clip+".dav_")
	}

	outputPath := filepath.Join(clipDirectory, clip+".mp4")

	fmt.Printf("converting %s into %s\n", inputPath, outputPath)

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

	data := ConvertingPageData{
		CameraName: camera,
		DayName:    day,
		VideoName:  video,
		HourName:   hour,
		ClipName:   clip,
	}

	err := tmpl.ExecuteTemplate(w, "converted.html", data)
	CheckTemplateError(w, err)
}
