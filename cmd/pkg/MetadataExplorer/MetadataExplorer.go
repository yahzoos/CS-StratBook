package MetadataExplorer

import (
	"encoding/json"
	"os"
	"strings"
)

// Metadata represents the structure of each entry in the JSON file
type Metadata struct {
	FileName     string `json:"file_name"`
	FilePath     string `json:"file_path"`
	ImagePath    string `json:"image_path"`
	NadeName     string `json:"nade_name"`
	Description  string `json:"description"`
	MapName      string `json:"map_name"`
	Side         string `json:"side"`
	NadeType     string `json:"nade_type"`
	SiteLocation string `json:"site"`
}

// Wrapper struct to correctly map the JSON file structure
type MetadataWrapper struct {
	Nades []Metadata `json:"nades"`
}

// LoadMetadata loads metadata from the fixed JSON file
func LoadMetadata(filePath string) ([]Metadata, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var wrapper MetadataWrapper
	err = json.Unmarshal(data, &wrapper)
	if err != nil {
		return nil, err
	}

	return wrapper.Nades, nil
}

func generateMaps(metadata []Metadata) []string {
	m := make(map[string]bool)
	var uniqueMaps []string
	for _, nades := range metadata {
		if !m[nades.MapName] {
			m[nades.MapName] = true
			uniqueMaps = append(uniqueMaps, nades.MapName)
		}
	}
	return uniqueMaps
}

// FilterOptions holds the selected user filters
type FilterOptions struct {
	MapPick  string
	T        bool
	CT       bool
	Smokes   bool
	Flashes  bool
	Molotovs bool
	HEs      bool
	ASite    bool
	BSite    bool
	MidSite  bool
}

var filters = FilterOptions{}

func FilterMetadata(metadata []Metadata, filters FilterOptions) []Metadata {
	var filtered []Metadata
	for _, nade := range metadata {
		if strings.ToLower(nade.MapName) != strings.ToLower(filters.MapPick) {
			continue
		}
		if (filters.T || filters.CT) &&
			!((filters.T && nade.Side == "T") || (filters.CT && nade.Side == "CT")) {
			continue
		}
		if (filters.Smokes || filters.Flashes || filters.Molotovs || filters.HEs) &&
			!((filters.Smokes && nade.NadeType == "smoke") ||
				(filters.Flashes && nade.NadeType == "flash") ||
				(filters.Molotovs && nade.NadeType == "molotov") ||
				(filters.Molotovs && nade.NadeType == "incendiary") ||
				(filters.HEs && nade.NadeType == "he")) {
			continue
		}
		if (filters.ASite || filters.BSite || filters.MidSite) &&
			!((filters.ASite && nade.SiteLocation == "A") ||
				(filters.BSite && nade.SiteLocation == "B") ||
				(filters.MidSite && nade.SiteLocation == "MID")) {
			continue
		}
		filtered = append(filtered, nade)
	}
	return filtered
}

type ReloadFunc func()
