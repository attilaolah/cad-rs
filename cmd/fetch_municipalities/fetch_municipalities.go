package main

import (
	"encoding/json"
	"log"
	"os"
	"sort"

	"github.com/attilaolah/ekat/municipalities"
)

func main() {
	all := []*municipalities.Municipality{}
	ms, errs := municipalities.FetchAll()
	go func() {
		for m := range ms {
			all = append(all, m)
		}
	}()

	for err := range errs {
		log.Fatalf("error fetching municipalities: %v", err)
	}

	sort.Slice(all, func(i, j int) bool { return all[i].ID < all[j].ID })
	json.NewEncoder(os.Stdout).Encode(&all)
}
