package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
)

type Clip struct {
	ClipName  string
	ClipID    string
	Debug     string
	Converted bool
}

type ClipListPageData struct {
	CameraName string
	DayName    string
	VideoName  string
	HourName   string
	Clips      []Clip
}

func ClipList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	camera := vars["camera"]
	day := vars["day"]
	video := vars["video"]
	hour := vars["hour"]

	directory := fmt.Sprintf("./camera-recordings/%s/%s/%s/%s", camera, day, video, hour)
	hour_directories, err := os.ReadDir(directory)
	CheckDirectoryError(w, err)

	var clips []Clip
	for _, clip := range hour_directories {
		fmt.Println(clip.Name())

		matched, _ := regexp.MatchString(".*\\.dav", clip.Name())
		if matched {
			// remove .dav or .dav_ from filename
			cleanClipName := strings.TrimSuffix(clip.Name(), ".dav")
			cleanClipName = strings.TrimSuffix(cleanClipName, ".dav_")

			clip := Clip{
				ClipName: cleanClipName,
				ClipID:   strings.ReplaceAll(cleanClipName, "[M][0@0][0]", ""),
			}

			// check other path if File exisis
			convertedPath := fmt.Sprintf("./converted/%s/%s/%s/%s.mp4", camera, day, hour, cleanClipName)
			clip.Debug = convertedPath
			_, err := os.Stat(convertedPath)
			if err == nil {
				clip.Converted = true
			} else {
				clip.Converted = false
			}
			clips = append(clips, clip)
		}
	}

	data := ClipListPageData{
		CameraName: camera,
		DayName:    day,
		VideoName:  video,
		HourName:   hour,
		Clips:      clips,
	}
	// fmt.Printf("%+v\n", data)
	err = tmpl.ExecuteTemplate(w, "clip_list.html", data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Fehler beim Ausführen des Templates: %v", err), http.StatusInternalServerError)
	}
}
