package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadCreatesDefaultWhenMissing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	want := Default()
	if cfg != want {
		t.Fatalf("got %+v, want %+v", cfg, want)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected config file to be created: %v", err)
	}
}

func TestSaveAndLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	in := Config{
		LastLevel: 4,
		MaxErrors: 3,
		MinWPM:    35,
	}
	if err := Save(path, in); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if got != in {
		t.Fatalf("got %+v, want %+v", got, in)
	}
}

func TestNormalizeBounds(t *testing.T) {
	cfg := Config{LastLevel: 0, MaxErrors: -2, MinWPM: 0}
	cfg.Normalize()
	if cfg.LastLevel != 1 || cfg.MaxErrors != 0 || cfg.MinWPM != 1 {
		t.Fatalf("unexpected normalized config: %+v", cfg)
	}
}
