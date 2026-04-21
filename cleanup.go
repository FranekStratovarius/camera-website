package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

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
