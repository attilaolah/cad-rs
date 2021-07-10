package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/attilaolah/ekat/scrapers"
)

var dst = flag.String("output_dir",
	filepath.Join(os.Getenv("BUILD_WORKSPACE_DIRECTORY"), "dist"),
	"Output directory (root) for scraped data.")

func main() {
	flag.Parse()

	ms, err := scrapers.ScrapeMunicipalities()
	if err != nil {
		log.Fatalf("failed to fetch municipalities: %v", err)
	}

	if err := scrapers.SaveMunicipalities(ms, *dst); err != nil {
		log.Fatalf("failed to save municipalities data: %v", err)
	}
}
