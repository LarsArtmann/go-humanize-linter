package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestComputeDelta_NoDelta(t *testing.T) {
	t.Parallel()

	entries := []baselineEntry{
		{Rule: "H001", File: "main.go", Line: 10},
		{Rule: "H002", File: "util.go", Line: 20},
	}

	delta := computeDelta(entries, entries)
	if delta.HasDelta() {
		t.Errorf("expected no delta, got added=%d removed=%d", len(delta.Added), len(delta.Removed))
	}
}

func TestComputeDelta_AddedFinding(t *testing.T) {
	t.Parallel()

	current := []baselineEntry{
		{Rule: "H001", File: "main.go", Line: 10},
		{Rule: "H002", File: "util.go", Line: 20},
	}
	baseline := []baselineEntry{
		{Rule: "H001", File: "main.go", Line: 10},
	}

	delta := computeDelta(current, baseline)
	if len(delta.Added) != 1 {
		t.Fatalf("expected 1 added, got %d", len(delta.Added))
	}

	if delta.Added[0].Rule != "H002" {
		t.Errorf("expected added H002, got %s", delta.Added[0].Rule)
	}

	if len(delta.Removed) != 0 {
		t.Errorf("expected 0 removed, got %d", len(delta.Removed))
	}
}

func TestComputeDelta_RemovedFinding(t *testing.T) {
	t.Parallel()

	current := []baselineEntry{
		{Rule: "H001", File: "main.go", Line: 10},
	}
	baseline := []baselineEntry{
		{Rule: "H001", File: "main.go", Line: 10},
		{Rule: "H003", File: "time.go", Line: 5},
	}

	delta := computeDelta(current, baseline)
	if len(delta.Removed) != 1 {
		t.Fatalf("expected 1 removed, got %d", len(delta.Removed))
	}

	if delta.Removed[0].Rule != "H003" {
		t.Errorf("expected removed H003, got %s", delta.Removed[0].Rule)
	}

	if len(delta.Added) != 0 {
		t.Errorf("expected 0 added, got %d", len(delta.Added))
	}
}

func TestLoadBaseline(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	data := baselineFile{
		Findings: []baselineEntry{
			{Rule: "H001", File: "main.go", Line: 10},
		},
	}

	raw, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	if err := os.WriteFile(path, raw, 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	entries, err := loadBaseline(path)
	if err != nil {
		t.Fatalf("loadBaseline: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	if entries[0].Rule != "H001" {
		t.Errorf("expected H001, got %s", entries[0].Rule)
	}
}

func TestLoadBaseline_MissingFile(t *testing.T) {
	t.Parallel()

	_, err := loadBaseline("/nonexistent/baseline.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
