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
		filepath.Join(os.Getenv("BUILD_WORKSPACE_DIRECTORY"), "dist"),
		"Output directory for scraped street data.")
	cache = flag.String("cache_dir",
		filepath.Join(os.Getenv("BUILD_WORKSPACE_DIRECTORY"), "dist", "street_search"),
		"Output directory for caching temporary scraped street search data.")
)

func main() {
	flag.Parse()

	m := int64(*mID)
	ss, errs := scrapers.ScrapeStreets(*cache, m)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for sr := range ss {
			if err := sr.Save(*cache, int64(*mID)); err != nil {
				log.Printf("error saving results to cache: %v", err)
			}
		}
	}()

	ok := true
	for err := range errs {
		log.Println(err)
		ok = false
	}

	wg.Wait()

	if !ok {
		os.Exit(1)
	}

	set, err := scrapers.MergeStreets(*cache, m)
	if err != nil {
		log.Fatalf("error merging scraped streets: %v", err)
	}

	if err := scrapers.SaveSettlements(set, *dst, m); err != nil {
		log.Fatalf("error saving scraped streets: %v", err)
	}
}
