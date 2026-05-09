package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"golang.org/x/image/draw"
)

func ThumbnailServer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	camera := vars["camera"]
	day := vars["day"]
	hour := vars["hour"]
	imageName := vars["image"]

	originalPath := filepath.Join("./camera-recordings", camera, day, "pic_001", hour, imageName+"[M][0@0][0].jpg")
	thumbDir := filepath.Join("./camera-recordings", camera, day, "thumbnails", hour)
	thumbPath := filepath.Join(thumbDir, imageName+".jpg")

	// Cache check: serve if the thumbnail is already generated on disk
	if _, err := os.Stat(thumbPath); err == nil {
		http.ServeFile(w, r, thumbPath)
		return
	}

	file, err := os.Open(originalPath)
	if err != nil {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		http.Error(w, "Error decoding image", http.StatusInternalServerError)
		return
	}

	bounds := img.Bounds()
	width := 400
	height := 300
	if bounds.Dx() > 0 {
		height = (bounds.Dy() * width) / bounds.Dx()
	}

	rect := image.Rect(0, 0, width, height)
	thumb := image.NewRGBA(rect)

	// Scale the image smoothly
	draw.BiLinear.Scale(thumb, rect, img, bounds, draw.Over, nil)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, thumb, &jpeg.Options{Quality: 80}); err != nil {
		http.Error(w, "Error encoding thumbnail", http.StatusInternalServerError)
		return
	}

	// Save generated thumbnail back to disk to serve fast next time
	if err := os.MkdirAll(thumbDir, 0755); err == nil {
		os.WriteFile(thumbPath, buf.Bytes(), 0644)
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))
	w.Write(buf.Bytes())
}
