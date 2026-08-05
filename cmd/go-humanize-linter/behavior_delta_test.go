package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/larsartmann/go-finding"
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

	if err := os.WriteFile(path, raw, 0o600); err != nil {
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

func TestSaveBaseline(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	report := finding.NewReportFromFindings(
		finding.ToolInfo{Name: "test"},
		[]finding.Finding{
			newFinding("H001", "a.go", 10),
			newFinding("H002", "b.go", 20),
		},
	)

	if err := saveBaseline(path, report); err != nil {
		t.Fatalf("saveBaseline: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read saved baseline: %v", err)
	}

	var decoded baselineFile
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal saved baseline: %v", err)
	}

	if len(decoded.Findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(decoded.Findings))
	}

	if decoded.Findings[0].Rule != "H001" || decoded.Findings[0].File != "a.go" || decoded.Findings[0].Line != 10 {
		t.Errorf("unexpected first entry: %+v", decoded.Findings[0])
	}
}

func TestReportToBaselineEntries(t *testing.T) {
	t.Parallel()

	report := finding.NewReportFromFindings(
		finding.ToolInfo{Name: "test"},
		[]finding.Finding{
			newFinding("H002", "b.go", 5),
			newFinding("H001", "a.go", 10),
		},
	)

	entries := reportToBaselineEntries(report)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	// Entries should be sorted by file, then line, then rule.
	if entries[0].Rule != "H001" || entries[0].File != "a.go" || entries[0].Line != 10 {
		t.Errorf("unexpected first sorted entry: %+v", entries[0])
	}
}

func TestPrintDelta_NoDelta(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	printDelta(&buf, behaviorDelta{})

	got := buf.String()
	if !strings.Contains(got, "No behavior delta") {
		t.Errorf("expected no-delta message, got: %q", got)
	}
}

func TestPrintDelta_AddedAndRemoved(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	printDelta(&buf, behaviorDelta{
		Added:   []baselineEntry{{Rule: "H001", File: "a.go", Line: 1}},
		Removed: []baselineEntry{{Rule: "H002", File: "b.go", Line: 2}},
	})

	got := buf.String()
	if !strings.Contains(got, "Added findings") {
		t.Errorf("expected added header, got: %q", got)
	}

	if !strings.Contains(got, "Removed findings") {
		t.Errorf("expected removed header, got: %q", got)
	}
}

func newFinding(rule, file string, line int) finding.Finding {
	return finding.NewBuilder(
		finding.RuleName(rule),
		finding.ToolName("test"),
		"test finding",
		finding.SeverityWarning,
		finding.Pos(finding.FilePath(file), line, 1),
	).
		WithConfidence(finding.ConfidenceHigh).
		MustBuild()
}
