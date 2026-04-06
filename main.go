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
	r.HandleFunc("/cameras/{camera}/{day}/{video}/{hour}/{clip}.mp4", ClipServer)
	r.HandleFunc("/convert/{camera}/{day}/{video}/{hour}/{clip}.mp4", ClipConverter)
	r.HandleFunc("/start-convert/{camera}/{day}/{video}/{hour}/{clip}.mp4", ConvertStart)

	// Statische Dateien aus dem "converted" Ordner bereitstellen, damit der Browser sie abspielen kann
	r.PathPrefix("/converted/").Handler(http.StripPrefix("/converted/", http.FileServer(http.Dir("./converted"))))

	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("test")
		pageData := PageData{
			PageTitle: "test",
			Test:      "blaaa",
		}

		tmpl := template.Must(template.ParseFiles("templates/basic.html"))
		tmpl.Execute(w, pageData)
	})

	r.HandleFunc("/clicked", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("clicked")
		pageData := PageData{
			PageTitle: "test",
			Test:      "blaaa",
		}

		tmpl := template.Must(template.ParseFiles("templates/test.html"))
		tmpl.Execute(w, pageData)
	})

	http.ListenAndServe(":80", r)
}
