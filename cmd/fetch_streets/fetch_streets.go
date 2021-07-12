package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/attilaolah/cad-rs/scrapers"
)

var (
	mID = flag.Int("municipality_id",
		80438, "Municipality ID to fetch streets for.")
	dst = flag.String("output_dir",
		filepath.Join(os.Getenv("BUILD_WORKSPACE_DIRECTORY"), "dist", "street_search"),
		"Output directory for scraped street search data.")
)

func main() {
	flag.Parse()

	ss, errs := scrapers.ScrapeStreets(*dst, int64(*mID))
	wg := sync.WaitGroup{}
	defer wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for sr := range ss {
			if err := sr.Save(*dst, int64(*mID)); err != nil {
				log.Printf("error saving results: %v", err)
			}
		}
	}()

	for err := range errs {
		log.Println(err)
	}
}
