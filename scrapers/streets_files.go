package scrapers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/attilaolah/cad-rs/proto"
)

// Save saves the results to disk.
func (sr *StreetSearchResults) Save(dir string, mID int64) error {
	subdir := filepath.Join(dir, strconv.FormatInt(mID, 10))
	if err := os.MkdirAll(subdir, DirPerm); err != nil {
		return fmt.Errorf("failed to create directory %q: %w", subdir, err)
	}

	return saveJSON(sr, subdir, asciil(sr.Query))
}

func MergeStreets(dir string, mID int64) ([]*pb.Settlement, error) {
	subdir := filepath.Join(dir, strconv.FormatInt(mID, 10))
	fs, err := ioutil.ReadDir(subdir)
	if err != nil {
		return nil, fmt.Errorf("failed to list contents of %q: %w", subdir, err)
	}

	sm := map[string]*pb.Settlement{}
	streets := map[string]map[string]time.Time{}
	ids := map[string]int64{}

	for _, fi := range fs {
		fn := filepath.Join(subdir, fi.Name())
		f, err := os.Open(fn)
		if err != nil {
			return nil, fmt.Errorf("failed to open %q: %w", fn, err)
		}

		r := StreetSearchResults{}
		err = json.NewDecoder(f).Decode(&r)
		cerr := f.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to decode file %q: %w", fn, err)
		}
		if cerr != nil {
			return nil, fmt.Errorf("failed to close file %q: %w", fn, err)
		}

		for _, s := range r.Results {
			parts := strings.SplitN(s.FullName, ", ", 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("failed to parse settlement from street full_name %q", s.FullName)
			}
			if _, ok := ids[s.FullName]; !ok {
				ids[s.FullName] = s.Id
			} else if ids[s.FullName] != s.Id {
				return nil, fmt.Errorf("mismatched IDs for street %q: %d vs %d", s.FullName, ids[s.FullName], s.Id)
			}

			set, str := parts[0], parts[1]
			if _, ok := sm[set]; !ok {
				sm[set] = &pb.Settlement{
					Name:      set,
					UpdatedAt: timestamppb.New(r.UpdatedAt),
				}
				streets[set] = map[string]time.Time{}
			} else if sm[set].UpdatedAt.AsTime().Before(r.UpdatedAt) {
				sm[set].UpdatedAt = timestamppb.New(r.UpdatedAt)
			}
			if streets[set][str].Before(r.UpdatedAt) {
				streets[set][str] = r.UpdatedAt
			}
		}
	}

	ss := []*pb.Settlement{}
	for _, s := range sm {
		ss = append(ss, s)
		for str, t := range streets[s.Name] {
			s.Streets = append(s.Streets, &pb.Street{
				Id:        ids[s.Name+", "+str],
				Name:      str,
				UpdatedAt: timestamppb.New(t),
			})
		}
		sort.Slice(s.Streets, func(i, j int) bool {
			return s.Streets[i].Name < s.Streets[j].Name
		})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Name < ss[j].Name
	})

	return ss, nil
}

func SaveSettlements(ss []*pb.Settlement, dir string, mID int64) error {
	subdir := filepath.Join(dir, "municipalities", strconv.FormatInt(mID, 10))
	// //municipalities/:id/settlements+streets.json
	if err := saveJSON(ss, subdir, "settlements+streets"); err != nil {
		return fmt.Errorf("failed to save municipalities/%d/settlements+streets: %w", mID, err)
	}

	return nil
}
