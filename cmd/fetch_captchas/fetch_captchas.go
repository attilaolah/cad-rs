package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/attilaolah/ekat/scrapers"
)

var municipalities = flag.String("municipalities",
	filepath.Join(os.Getenv("BUILD_WORKSPACE_DIRECTORY"), "data", "municipalities.json"),
	"JSON file containing municipalities.")

func main() {
	flag.Parse()

	ctx := context.Background()

	// Ctrl+C to cancel the context.
	ctx, cancel := context.WithCancel(ctx)
	cancelled := false

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	defer func() {
		signal.Stop(sig)
		cancel()
	}()

	go func() {
		select {
		case <-sig:
			cancelled = true
			cancel()
		case <-ctx.Done():
		}
	}()

	n := 0
	cs, errs := scrapers.Scrape4Captchas(ctx, *municipalities, 2)

	for {
		select {
		case c := <-cs:
			json.NewEncoder(os.Stdout).Encode(c)
			n += 1
		case err := <-errs:
			log.Printf("error: %v", err)
		case <-ctx.Done():
			if cancelled {
				// Ctrl+C caught, quit nicely.
				fmt.Printf("\rQUIT: %d captchas fetched\n", n)
				return
			}
			log.Fatalf("quitting: %v", ctx.Err())
		}
	}
}
