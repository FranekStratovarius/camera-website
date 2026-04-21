package main

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"

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

func passwordAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// r.BasicAuth() liest die Header aus. Den Usernamen (_) ignorieren wir.
		_, password, ok := r.BasicAuth()

		// Sicherheits-Feature: subtle.ConstantTimeCompare verhindert sogenannte "Timing-Attacken".
		// Es vergleicht die Strings immer in der exakt gleichen Zeit, egal ob das Passwort
		// beim ersten oder letzten Zeichen falsch ist.
		if !ok || subtle.ConstantTimeCompare([]byte(password), []byte(os.Getenv("PASSWORD"))) != 1 {
			// Trigger für das Browser-Pop-up
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted Area", charset="UTF-8"`)
			http.Error(w, "Zugriff verweigert", http.StatusUnauthorized)
			return
		}

		// Passwort ist korrekt -> Weiter zur eigentlichen Funktion
		next.ServeHTTP(w, r)
	}
}

func cleanupRoutine() {
	cleanupOldFolders("camera-recordings/Kamera_Hof", 3)

	ticker := time.NewTicker(24 * time.Hour)
	for range ticker.C {
		cleanupOldFolders("camera-recordings/Kamera_Hof", 3)
	}
}

func cleanupOldFolders(basePath string, keep int) {
	entries, err := os.ReadDir(basePath)
	if err != nil {
		fmt.Printf("Fehler beim Lesen des Verzeichnisses %s: %v\n", basePath, err)
		return
	}

	var folders []os.DirEntry
	for _, entry := range entries {
		// ignore folderst starting with $
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), "$") {
			folders = append(folders, entry)
		}
	}

	// sort folders chronologically by their name (needs format "YYYY-MM-DD")
	sort.Slice(folders, func(i, j int) bool {
		return folders[i].Name() > folders[j].Name()
	})

	// delete all older folders
	if len(folders) > keep {
		for _, folder := range folders[keep:] {
			dirToRemove := filepath.Join(basePath, folder.Name())
			if err := os.RemoveAll(dirToRemove); err != nil {
				fmt.Printf("Fehler beim Löschen des Ordners %s: %v\n", dirToRemove, err)
			} else {
				fmt.Printf("Alter Aufnahme-Ordner automatisch gelöscht: %s\n", dirToRemove)
			}
		}
	}
}

func main() {
	go cleanupRoutine()

	r := mux.NewRouter()

	r.HandleFunc("/", passwordAuth(CameraList))
	r.HandleFunc("/cameras/{camera}", passwordAuth(DayList))
	r.HandleFunc("/cameras/{camera}/{day}", passwordAuth(HourList))
	r.HandleFunc("/cameras/{camera}/{day}/{hour}", passwordAuth(ClipList))
	r.HandleFunc("/convert/{camera}/{day}/{hour}/{clip}.mp4", passwordAuth(ClipConverter))
	r.HandleFunc("/converted/{camera}/{day}/{hour}/{clip}.mp4", passwordAuth(ClipServer))
	r.HandleFunc("/images/{camera}/{day}/{hour}/{image}.jpg", passwordAuth(ImageServer))
	r.HandleFunc("/start-convert/{camera}/{day}/{hour}/{clip}.mp4", passwordAuth(ConvertStart))

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.ListenAndServe(":80", r)
}
