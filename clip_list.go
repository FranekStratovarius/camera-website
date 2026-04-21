package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
)

type Clip struct {
	ClipName  string
	ClipID    string
	Debug     string
	Converted bool
	Time      string
}

type Image struct {
	ImageName string
	Time      string
}

type ClipListPageData struct {
	CameraName       string
	DayName          string
	DayNameFormatted string
	HourName         string
	Clips            []Clip
	Images           []Image
}

func ClipList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	camera := vars["camera"]
	day := vars["day"]
	hour := vars["hour"]

	directory := filepath.Join("./camera-recordings", camera, day, "video_001", hour)
	hour_directories, err := os.ReadDir(directory)
	CheckDirectoryError(w, err)

	var clips []Clip
	for _, clip := range hour_directories {

		matched, _ := regexp.MatchString(".*\\.dav", clip.Name())
		if matched {
			// remove .dav or .dav_ from filename
			cleanClipName := strings.TrimSuffix(clip.Name(), ".dav")
			cleanClipName = strings.TrimSuffix(cleanClipName, ".dav_")

			clip := Clip{
				ClipName: cleanClipName,
				ClipID:   strings.ReplaceAll(cleanClipName, "[M][0@0][0]", ""),
			}

			times := strings.Split(strings.ReplaceAll(cleanClipName, "[M][0@0][0]", ""), "-")
			startTimeParts := strings.Split(times[0], ".")
			endTimeParts := strings.Split(times[1], ".")
			clip.Time = fmt.Sprintf("%s:%s:%s - %s:%s:%s", startTimeParts[0], startTimeParts[1], startTimeParts[2], endTimeParts[0], endTimeParts[1], endTimeParts[2])

			// check other path if File exisis
			convertedPath := filepath.Join(directory, cleanClipName+".mp4")
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

	directory = filepath.Join("./camera-recordings", camera, day, "pic_001", hour)
	image_directory, err := os.ReadDir(directory)
	CheckDirectoryError(w, err)

	var images []Image
	for _, image := range image_directory {
		imageCleanName := strings.TrimSuffix(image.Name(), ".jpg")
		imageCleanName = strings.ReplaceAll(imageCleanName, "[M][0@0][0]", "")
		imageTimeParts := strings.Split(imageCleanName, ".")
		time := fmt.Sprintf("%s:%s:%s", imageTimeParts[0], imageTimeParts[1], imageTimeParts[2])

		image := Image{
			ImageName: imageCleanName,
			Time:      time,
		}
		images = append(images, image)
	}

	dayNameParts := strings.Split(day, "-")
	dayNameFormatted := fmt.Sprintf("%s.%s.%s", dayNameParts[2], dayNameParts[1], dayNameParts[0])

	data := ClipListPageData{
		CameraName:       camera,
		DayName:          day,
		DayNameFormatted: dayNameFormatted,
		HourName:         hour,
		Clips:            clips,
		Images:           images,
	}

	err = tmpl.ExecuteTemplate(w, "clip_list.html", data)
	CheckTemplateError(w, err)
}
