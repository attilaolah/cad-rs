package scrapers

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// Save saves the results to disk.
func (sr *StreetSearchResults) Save(dir string, mID int64) error {
	subdir := filepath.Join(dir, strconv.FormatInt(mID, 10))
	if err := os.MkdirAll(subdir, DirPerm); err != nil {
		return fmt.Errorf("failed to create directory %q: %w", subdir, err)
	}

	return saveJSON(sr, subdir, asciil(sr.Query))
}
