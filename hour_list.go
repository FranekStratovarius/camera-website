package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

type Hour struct {
	HourName       string
	VideoDirectory string
}

type HourListPageData struct {
	CameraName       string
	DayName          string
	DayNameFormatted string
	Hours            []Hour
}

func HourList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	camera := vars["camera"]
	day := vars["day"]

	directory := fmt.Sprintf("./camera-recordings/%s/%s", camera, day)
	video_directories, err := os.ReadDir(directory)
	CheckDirectoryError(w, err)

	var hours []Hour
	for _, video := range video_directories {
		directory := fmt.Sprintf("./camera-recordings/%s/%s/%s", camera, day, video.Name())
		fileInfo, err := os.Stat(directory)
		CheckDirectoryError(w, err)

		if fileInfo.IsDir() {
			video_directories, err := os.ReadDir(directory)
			CheckDirectoryError(w, err)

			for _, hour := range video_directories {
				directory := fmt.Sprintf("./camera-recordings/%s/%s/%s/%s", camera, day, video.Name(), hour.Name())
				fileInfo, err := os.Stat(directory)
				CheckDirectoryError(w, err)

				if fileInfo.IsDir() {
					hours = append(hours, Hour{
						HourName:       hour.Name(),
						VideoDirectory: video.Name(),
					})
				}
			}
		}
	}

	dayNameParts := strings.Split(day, "-")
	dayNameFormatted := fmt.Sprintf("%s.%s.%s", dayNameParts[2], dayNameParts[1], dayNameParts[0])

	data := HourListPageData{
		CameraName:       camera,
		DayName:          day,
		DayNameFormatted: dayNameFormatted,
		Hours:            hours,
	}

	err = tmpl.ExecuteTemplate(w, "hour_list.html", data)
	CheckTemplateError(w, err)
}
