package main

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/larsartmann/go-finding"
)

// baselineEntry is a single finding recorded in a baseline file. Only the
// stable identity fields are stored — message text changes between versions
// and should not trigger a delta.
type baselineEntry struct {
	Rule string `json:"rule"`
	File string `json:"file"`
	Line int    `json:"line"`
}

// baselineFile is the JSON structure of a --behavior-delta baseline.
type baselineFile struct {
	Findings []baselineEntry `json:"findings"`
}

// baselineFileMode is the file mode used when writing baseline JSON files.
// It is restrictive because baselines may be checked into source control.
const baselineFileMode = 0o600

// findingKey is the stable identity tuple used for delta comparison.
type findingKey struct {
	Rule string
	File string
	Line int
}

// behaviorDelta represents the difference between the current run and a
// saved baseline. Added findings are potential new false positives; removed
// findings are potential missed detections.
type behaviorDelta struct {
	Added   []baselineEntry
	Removed []baselineEntry
}

// HasDelta reports whether any findings were added or removed.
func (d behaviorDelta) HasDelta() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0
}

// loadBaseline reads a JSON baseline file and returns its findings.
func loadBaseline(path string) ([]baselineEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read baseline: %w", err)
	}

	var bf baselineFile
	if err := json.Unmarshal(data, &bf); err != nil {
		return nil, fmt.Errorf("parse baseline JSON: %w", err)
	}

	return bf.Findings, nil
}

// saveBaseline writes a baseline file from a report's findings.
func saveBaseline(path string, report *finding.Report) error {
	entries := reportToBaselineEntries(report)

	bf := baselineFile{Findings: entries}

	data, err := json.Marshal(bf, jsontext.WithIndentPrefix(""), jsontext.WithIndent("  "))
	if err != nil {
		return fmt.Errorf("marshal baseline: %w", err)
	}

	data = append(data, '\n')

	if err := os.WriteFile(path, data, baselineFileMode); err != nil {
		return fmt.Errorf("write baseline: %w", err)
	}

	return nil
}

// computeDelta compares current findings against a baseline and returns the
// difference. Comparison is by (rule, file, line) tuple — message text is
// not compared because it may change between versions without indicating a
// behavioral regression.
func computeDelta(current []baselineEntry, baseline []baselineEntry) behaviorDelta {
	currentSet := make(map[findingKey]baselineEntry, len(current))
	for _, entry := range current {
		currentSet[findingKey(entry)] = entry
	}

	baselineSet := make(map[findingKey]baselineEntry, len(baseline))
	for _, entry := range baseline {
		baselineSet[findingKey(entry)] = entry
	}

	var delta behaviorDelta

	for key, entry := range currentSet {
		if _, ok := baselineSet[key]; !ok {
			delta.Added = append(delta.Added, entry)
		}
	}

	for key, entry := range baselineSet {
		if _, ok := currentSet[key]; !ok {
			delta.Removed = append(delta.Removed, entry)
		}
	}

	sortEntries(delta.Added)
	sortEntries(delta.Removed)

	return delta
}

// printDelta writes the delta to the given writer in a human-readable format.
func printDelta(writer io.Writer, delta behaviorDelta) {
	if len(delta.Added) > 0 {
		fmt.Fprintf(writer, "Added findings (%d — potential new false positives):\n", len(delta.Added))

		for _, entry := range delta.Added {
			fmt.Fprintf(writer, "  + %s %s:%d\n", entry.Rule, entry.File, entry.Line)
		}
	}

	if len(delta.Removed) > 0 {
		fmt.Fprintf(writer, "Removed findings (%d — potential missed detections):\n", len(delta.Removed))

		for _, entry := range delta.Removed {
			fmt.Fprintf(writer, "  - %s %s:%d\n", entry.Rule, entry.File, entry.Line)
		}
	}

	if !delta.HasDelta() {
		fmt.Fprintln(writer, "No behavior delta — findings match baseline exactly.")
	}
}

// reportToBaselineEntries converts a report's findings to baseline entries.
func reportToBaselineEntries(report *finding.Report) []baselineEntry {
	var entries []baselineEntry

	for finding := range report.All() {
		entries = append(entries, baselineEntry{
			Rule: string(finding.Rule),
			File: string(finding.Position.File),
			Line: finding.Position.Line,
		})
	}

	sortEntries(entries)

	return entries
}

func sortEntries(entries []baselineEntry) {
	sort.Slice(entries, func(left, right int) bool {
		if entries[left].File != entries[right].File {
			return entries[left].File < entries[right].File
		}

		if entries[left].Line != entries[right].Line {
			return entries[left].Line < entries[right].Line
		}

		return entries[left].Rule < entries[right].Rule
	})
}
