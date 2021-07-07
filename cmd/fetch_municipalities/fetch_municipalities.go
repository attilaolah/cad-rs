package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/attilaolah/ekat/municipalities"
)

var (
	output = flag.String("output",
		func() string {
			workspace := os.Getenv("BUILD_WORKSPACE_DIRECTORY")
			if workspace != "" {
				return filepath.Join(workspace, "data", "municipalities.json")
			}
			return "-"
		}(),
		"Where to write the output file (- means stdout).")
)

func main() {
	flag.Parse()

	ms, err := municipalities.FetchAll()
	if err != nil {
		log.Fatalf("error fetching municipalities: %v", err)
	}
	data, err := json.MarshalIndent(ms, "", "  ")
	if err != nil {
		log.Fatalf("error encoding municipalities: %v", err)
	}

	var out *os.File = os.Stdout
	if *output != "-" {
		f, err := os.Create(*output)
		if err != nil {
			log.Fatalf("error creating output file: %v", err)
		}
		out = f
		defer func() {
			if err = out.Close(); err != nil {
				log.Fatalf("error closing output file: %v", err)
			}
		}()
	}

	data = append(data, '\n')
	if _, err = out.Write(data); err != nil {
		log.Fatalf("error writing to output file: %v", err)
	}
}
