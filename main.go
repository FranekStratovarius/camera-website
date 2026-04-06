package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

var tmpl *template.Template

type PageData struct {
	PageTitle string
	Test      string
}

func CheckDirectoryError(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Fehler beim Lesen des Verzeichnisses: %v", err),
			http.StatusInternalServerError,
		)
		return
	}
}

func CheckTemplateError(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Fehler beim Ausführen des Templates: %v", err),
			http.StatusInternalServerError,
		)
		return
	}
}

func init() {
	// Lädt einmalig beim Programmstart alle .html Dateien aus dem templates-Ordner
	tmpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", CameraList)
	r.HandleFunc("/cameras/{camera}", DayList)
	r.HandleFunc("/cameras/{camera}/{day}", HourList)
	r.HandleFunc("/cameras/{camera}/{day}/{video}/{hour}", ClipList)
	r.HandleFunc("/convert/{camera}/{day}/{video}/{hour}/{clip}.mp4", ClipConverter)
	r.HandleFunc("/converted/{camera}/{day}/{video}/{hour}/{clip}.mp4", ClipServer)
	r.HandleFunc("/start-convert/{camera}/{day}/{video}/{hour}/{clip}.mp4", ConvertStart)

	http.ListenAndServe(":80", r)
}
