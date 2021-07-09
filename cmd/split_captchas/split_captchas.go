package main

import (
	"flag"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/attilaolah/ekat/labeller"
)

var datadir = flag.String("data_dir",
	filepath.Join(os.Getenv("BUILD_WORKSPACE_DIRECTORY"), "data", "captchas"),
	"Directory containing gaptcha files.")

func main() {
	flag.Parse()

	imgs, errs := labeller.Split4Captchas(*datadir)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var img image.Image
		select {
		case img = <-imgs:
		default:
			w.Header().Set("content-type", "text/plain")
			http.Error(w, "image not found", http.StatusNotFound)
			return
		}

		w.Header().Set("content-type", "image/png")
		if err := png.Encode(w, img); err != nil {
			log.Printf("error writing response: %v", err)
		}
	})
	go log.Fatal(http.ListenAndServe(":8080", nil))

	for err := range errs {
		log.Printf("error: %v", err)
	}
}
