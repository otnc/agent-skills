package main

import (
	"encoding/json"
	"log"
	"os"
)

const (
	configVersion = 1
	stateVersion  = 1
)

// Config is sources.json: where to discover skill repositories and which ones to skip. Discovery is automatic; the file only tunes it.
type Config struct {
	Version int      `json:"version"`
	Owner   string   `json:"owner"`
	Prefix  string   `json:"prefix"`
	Exclude []string `json:"exclude"` // "owner/repo" or bare repo names to never sync
}

// LicenseInfo records the license of a source repository, so the generated docs can link it next to the repository.
type LicenseInfo struct {
	SpdxID string `json:"spdxId"` // SPDX identifier, e.g. "MIT"; "NOASSERTION" when GitHub cannot classify the file
	URL    string `json:"url"`    // link to the license file in the source repository
}

// SourceState is the per-source entry of sync-state.json.
type SourceState struct {
	Tag      string       `json:"tag"`  // last synced release tag
	Dirs     []string     `json:"dirs"` // skill directories under skills/ owned by this source
	License  *LicenseInfo `json:"license,omitempty"`
	SyncedAt string       `json:"syncedAt"`
}

// State is sync-state.json: what has been synced so far.
type State struct {
	Version int                     `json:"version"`
	Sources map[string]*SourceState `json:"sources"`
}

func mustLoadConfig(path string) *Config {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("read %s: %v", path, err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("parse %s: %v", path, err)
	}
	if cfg.Owner == "" {
		log.Fatalf("%s: owner is required", path)
	}
	if cfg.Prefix == "" {
		cfg.Prefix = "skills_"
	}
	return &cfg
}

func mustLoadState(path string) *State {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &State{Version: stateVersion, Sources: map[string]*SourceState{}}
	}
	if err != nil {
		log.Fatalf("read %s: %v", path, err)
	}
	var st State
	if err := json.Unmarshal(data, &st); err != nil {
		log.Fatalf("parse %s: %v", path, err)
	}
	if st.Sources == nil {
		st.Sources = map[string]*SourceState{}
	}
	return &st
}

func saveState(path string, st *State) {
	st.Version = stateVersion
	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		log.Fatalf("encode %s: %v", path, err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		log.Fatalf("write %s: %v", path, err)
	}
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func logf(format string, args ...any) {
	log.Printf(format, args...)
}
