package main

import (
	"fmt"
	"net/http"
	"os"
)

type Camera struct {
	CameraName string
}

type CameraListPageData struct {
	Cameras []Camera
}

func CameraList(w http.ResponseWriter, r *http.Request) {
	directory := "./camera-recordings"
	cameras_directories, err := os.ReadDir(directory)
	CheckDirectoryError(w, err)

	var cameras []Camera
	for _, camera := range cameras_directories {
		directory := fmt.Sprintf("./camera-recordings/%s", camera.Name())
		fileInfo, err := os.Stat(directory)
		CheckDirectoryError(w, err)

		if fileInfo.IsDir() {
			cameras = append(cameras, Camera{CameraName: camera.Name()})
		}
	}

	data := CameraListPageData{
		Cameras: cameras,
	}
	err = tmpl.ExecuteTemplate(w, "camera_list.html", data)
	CheckTemplateError(w, err)
}
