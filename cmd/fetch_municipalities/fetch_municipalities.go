package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/attilaolah/ekat/municipalities"
)

func main() {
	ms, err := municipalities.FetchMunicipalities()
	if err != nil {
		log.Fatalf("error fetching municipalities: %v", err)
		return
	}
	if err := json.NewEncoder(os.Stdout).Encode(&ms); err != nil {
		log.Fatalf("error encoding municipalities: %v", err)
	}
}
