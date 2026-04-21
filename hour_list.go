package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

type Hour struct {
	HourName string
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

	directory := fmt.Sprintf("./camera-recordings/%s/%s/video_001", camera, day)
	hours_directories, err := os.ReadDir(directory)
	CheckDirectoryError(w, err)

	var hours []Hour
	for _, hour := range hours_directories {
		directory := fmt.Sprintf("./camera-recordings/%s/%s/video_001/%s", camera, day, hour.Name())
		fileInfo, err := os.Stat(directory)
		CheckDirectoryError(w, err)

		if fileInfo.IsDir() {
			hours = append(hours, Hour{
				HourName: hour.Name(),
			})
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
