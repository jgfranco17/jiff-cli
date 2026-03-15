package diffs

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/fatih/color"
)

type ComparisonResult struct {
	Added   []string
	Removed []string
	Changed []string
}

func (r ComparisonResult) IsEmpty() bool {
	return len(r.Added) == 0 && len(r.Removed) == 0 && len(r.Changed) == 0
}

func (r ComparisonResult) Total() int {
	return len(r.Added) + len(r.Removed) + len(r.Changed)
}

func (r ComparisonResult) Render(w io.Writer) {
	fmt.Fprintf(w, "%d differences found:\n", r.Total())
	redSprintf := color.New(color.FgRed).SprintfFunc()
	yellowSprintf := color.New(color.FgYellow).SprintfFunc()
	greenSprintf := color.New(color.FgGreen).SprintfFunc()

	renderList(w, r.Removed, "- ", redSprintf)
	renderList(w, r.Changed, "~ ", yellowSprintf)
	renderList(w, r.Added, "+ ", greenSprintf)
}

// CompareJSON returns a ComparisonResult containing human-readable difference
// strings between two JSON objects. Keys present only in source are reported as
// removed, keys present only in target as added, and mismatched values as
// changed.
func CompareJSON(source, target io.Reader) (ComparisonResult, error) {
	var sourceJSON, targetJSON map[string]any
	if err := json.NewDecoder(source).Decode(&sourceJSON); err != nil {
		return ComparisonResult{}, err
	}
	if err := json.NewDecoder(target).Decode(&targetJSON); err != nil {
		return ComparisonResult{}, err
	}
	var comparison ComparisonResult

	for k, sv := range sourceJSON {
		tv, ok := targetJSON[k]
		if !ok {
			comparison.Removed = append(comparison.Removed, fmt.Sprintf("%q: %v", k, sv))
			continue
		}
		if fmt.Sprintf("%v", sv) != fmt.Sprintf("%v", tv) {
			comparison.Changed = append(comparison.Changed, fmt.Sprintf("%q: %v -> %v", k, sv, tv))
		}
	}
	for k, tv := range targetJSON {
		if _, ok := sourceJSON[k]; !ok {
			comparison.Added = append(comparison.Added, fmt.Sprintf("%q: %v", k, tv))
		}
	}

	return comparison, nil
}

func renderList(w io.Writer, items []string, prefix string, colorFunc func(format string, a ...any) string) {
	for _, item := range items {
		fmt.Fprintln(w, colorFunc(prefix+item))
	}
}
