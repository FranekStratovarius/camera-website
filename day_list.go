package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Day struct {
	DayName string
}

type DayListPageData struct {
	CameraName string
	Days       []Day
}

func DayList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	camera := vars["camera"]

	directory := fmt.Sprintf("./camera-recordings/%s", camera)
	days_directories, err := os.ReadDir(directory)
	CheckDirectoryError(w, err)

	var days []Day
	for _, day := range days_directories {
		fmt.Println(day.Name())
		directory := fmt.Sprintf("./camera-recordings/%s/%s", camera, day.Name())
		fileInfo, err := os.Stat(directory)
		CheckDirectoryError(w, err)

		if fileInfo.IsDir() {
			days = append(days, Day{DayName: day.Name()})
		}
	}

	data := DayListPageData{
		CameraName: camera,
		Days:       days,
	}
	err = tmpl.ExecuteTemplate(w, "day_list.html", data)
	CheckTemplateError(w, err)
}
