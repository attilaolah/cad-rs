package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/attilaolah/ekat/proto"
	"github.com/attilaolah/ekat/scrapers"
)

var (
	municipalities = flag.String("municipalities",
		filepath.Join(os.Getenv("BUILD_WORKSPACE_DIRECTORY"), "data", "municipalities.json"),
		"JSON file containing municipalities.")
	dst = flag.String("output_dir",
		filepath.Join(os.Getenv("BUILD_WORKSPACE_DIRECTORY"), "data", "captchas"),
		"Output directory for scraped metadata and images.")
)

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
	sdir := filepath.Join(*dst, "samples")
	cs, errs := scrapers.Scrape4Captchas(ctx, *municipalities, 2)

	for {
		select {
		case c := <-cs:
			for _, s := range c.Samples {
				fn := filepath.Join(sdir, fmt.Sprintf("%s.jpg", s.Sha1))
				f, err := os.Create(fn)
				if err != nil {
					log.Printf("error creating sample file: %v", err)
					continue
				}
				if _, err = f.Write(s.Data); err != nil {
					log.Printf("error writing sample data: %v", err)
					f.Close()
					continue
				}
				if err = f.Close(); err != nil {
					log.Printf("error closing sample file: %v", err)
					continue
				}
				// Omit JSON binary image data.
				s.Data = nil
			}

			tmp, err := save(c)
			if err != nil {
				log.Printf("error saving captcha: %v", err)
				break
			}
			if err = rename(tmp); err != nil {
				log.Printf("error renaming captcha: %v", err)
				break
			}
			n += 1
			fmt.Printf("\rSAVE: %q [%d]", c.Id, n)
			os.Stdout.Sync()
		case err := <-errs:
			log.Printf("error fetching captcha: %v", err)
		case <-ctx.Done():
			if cancelled {
				// Ctrl+C caught, quit nicely.
				fmt.Printf("\nQUIT: %d captchas fetched\n", n)
				return
			}
			log.Fatalf("quitting: %v", ctx.Err())
		}
	}
}

// Save the captcha metadata to a file, returning the filename.
// The returned filename should be renamed (atomically) to its final name.
func save(c *proto.Captcha) (string, error) {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error encoding %q: %w", c.Id, err)
	}
	data = append(data, '\n')

	p := fmt.Sprintf("%s.*.json", c.Id)
	tmp, err := ioutil.TempFile(*dst, p)
	if err != nil {
		return "", fmt.Errorf("error creating file with pattern %q: %w", p, err)
	}

	defer func() {
		cerr := tmp.Close()
		if err == nil {
			err = fmt.Errorf("error closing file: %w", cerr)
		}
	}()

	if _, err = tmp.Write(data); err != nil {
		return "", fmt.Errorf("error writing %q: %w", tmp.Name(), err)
	}

	return tmp.Name(), err
}

// Rename a temporary file to its final name.
func rename(tmp string) error {
	dir := filepath.Dir(tmp)
	dst := fmt.Sprintf("%s.json", filepath.Base(tmp)[:36])
	return os.Rename(tmp, filepath.Join(dir, dst))
}
