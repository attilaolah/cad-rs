package scrapers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	pb "github.com/attilaolah/cad-rs/proto"
)

// DirPerm encodes new directory permissions.
const DirPerm = 0o755

// ScalarMunicipality is a Municipality with only scalar fields.
// Non-scalar (i.e. message) fields are turned into references (i.e. IDs).
type ScalarMunicipality struct {
	pb.Municipality

	CadastralMunicipalities []int64 `json:"cadastral_municipalities"`
}

// SaveMunicipalities stores municipality data in the expected directory layout.
func SaveMunicipalities(ms []*pb.Municipality, dir string) error {
	{
		data := make([]ScalarMunicipality, len(ms))
		for i, m := range ms {
			data[i].Municipality = *m
			data[i].Municipality.CadastralMunicipalities = nil
			data[i].CadastralMunicipalities = make([]int64, len(m.CadastralMunicipalities))
			for j, cm := range m.CadastralMunicipalities {
				data[i].CadastralMunicipalities[j] = cm.Id
			}
		}

		// /municipalities.json
		if err := saveJSON(data, dir, "municipalities"); err != nil {
			return fmt.Errorf("failed to save municipalities: %w", err)
		}
	}

	subms := filepath.Join(dir, "municipalities")
	if err := os.MkdirAll(subms, DirPerm); err != nil {
		return fmt.Errorf("failed to create directory %q: %w", subms, err)
	}

	{
		ids := make([]int64, len(ms))
		for i, m := range ms {
			ids[i] = m.Id
		}

		// /municipalities/ids.json
		if err := saveJSON(ids, subms, "ids"); err != nil {
			return fmt.Errorf("failed to save municipalities/ids: %w", err)
		}
	}

	for _, m := range ms {
		// /municipalities/:id.json
		if err := saveJSON(m, subms, fmt.Sprintf("%d", m.Id)); err != nil {
			return fmt.Errorf("failed to save municipalities/%d: %w", m.Id, err)
		}

		subm := filepath.Join(subms, fmt.Sprintf("%d", m.Id))
		if err := os.MkdirAll(subm, DirPerm); err != nil {
			return fmt.Errorf("failed to create directory %q: %w", subm, err)
		}

		// /municipalities/:id/cadastral_municipalities.json
		if err := saveJSON(m.CadastralMunicipalities, subm, "cadastral_municipalities"); err != nil {
			return fmt.Errorf("failed to save municipalities/%d/cadastral_municipalities: %w", m.Id, err)
		}

		subcm := filepath.Join(subm, "cadastral_municipalities")
		if err := os.MkdirAll(subcm, DirPerm); err != nil {
			return fmt.Errorf("failed to create directory %q: %w", subcm, err)
		}

		ids := make([]int64, len(m.CadastralMunicipalities))
		for i, cm := range m.CadastralMunicipalities {
			ids[i] = cm.Id
		}

		// /municipalities/:id/cadastral_municipalities/ids.json
		if err := saveJSON(ids, subcm, "ids"); err != nil {
			return fmt.Errorf("failed to save municipalities/%d/cadastral_municipalities/ids: %w", m.Id, err)
		}

		for _, cm := range m.CadastralMunicipalities {
			// /municipalities/:id/cadastral_municipalities/:id.json
			if err := saveJSON(cm, subcm, fmt.Sprintf("%d", cm.Id)); err != nil {
				return fmt.Errorf("failed to save municipalities/%d/cadastral_municipalities/%d: %w", m.Id, cm.Id, err)
			}
		}
	}

	// //municipalities+cadastral_municipalities.json
	if err := saveJSON(ms, dir, "municipalities+cadastral_municipalities"); err != nil {
		return fmt.Errorf("failed to save municipalities+cadastral_municipalities: %w", err)
	}

	return nil
}

// Saves JSON data as dir/fn.min.json and dir/fn.json.
func saveJSON(data interface{}, dir, fn string) error {
	raw, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode data: %w", err)
	}
	if err := saveFile(append(raw, '\n'), filepath.Join(dir, fn)+".json"); err != nil {
		return fmt.Errorf("failed to save data: %w", err)
	}

	if raw, err = json.Marshal(data); err != nil {
		return fmt.Errorf("failed to encode data (minified): %w", err)
	}
	if err := saveFile(raw, filepath.Join(dir, fn)+".min.json"); err != nil {
		return fmt.Errorf("failed to save data (minified): %w", err)
	}

	return nil
}

// Saves raw data into a file, handling create/write/close.
func saveFile(data []byte, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %q: %w", path, err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			err = fmt.Errorf("failed to close file %q: %w", path, cerr)
		}
	}()

	if _, err = f.Write(data); err != nil {
		return fmt.Errorf("failed to write data to file %q: %w", path, err)
	}

	return err
}
