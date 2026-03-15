package diffs

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompareJSON(t *testing.T) {
	tests := []struct {
		name        string
		source      string
		target      string
		wantErr     bool
		wantAdded   []string
		wantRemoved []string
		wantChanged []string
	}{
		{
			name:        "identical objects produce empty result",
			source:      `{"a":"1","b":"2"}`,
			target:      `{"a":"1","b":"2"}`,
			wantAdded:   nil,
			wantRemoved: nil,
			wantChanged: nil,
		},
		{
			name:        "key added in target",
			source:      `{"a":"1"}`,
			target:      `{"a":"1","b":"2"}`,
			wantAdded:   []string{`"b": 2`},
			wantRemoved: nil,
			wantChanged: nil,
		},
		{
			name:        "key removed from source",
			source:      `{"a":"1","b":"2"}`,
			target:      `{"a":"1"}`,
			wantAdded:   nil,
			wantRemoved: []string{`"b": 2`},
			wantChanged: nil,
		},
		{
			name:        "value changed",
			source:      `{"a":"old"}`,
			target:      `{"a":"new"}`,
			wantAdded:   nil,
			wantRemoved: nil,
			wantChanged: []string{`"a": old -> new`},
		},
		{
			name:        "simultaneous add remove and change",
			source:      `{"keep":"x","remove":"y","change":"old"}`,
			target:      `{"keep":"x","add":"z","change":"new"}`,
			wantAdded:   []string{`"add": z`},
			wantRemoved: []string{`"remove": y`},
			wantChanged: []string{`"change": old -> new`},
		},
		{
			name:        "nested map order does not affect comparison",
			source:      `{"nest": {"a":"1","b":[1, 2, 3]}}`,
			target:      `{"nest": {"b":[1, 2, 3],"a":"1"}}`,
			wantAdded:   nil,
			wantRemoved: nil,
			wantChanged: nil,
		},
		{
			name:    "invalid JSON in source returns error",
			source:  `not-json`,
			target:  `{"a":"1"}`,
			wantErr: true,
		},
		{
			name:    "invalid JSON in target returns error",
			source:  `{"a":"1"}`,
			target:  `{broken`,
			wantErr: true,
		},
		{
			name:        "empty objects produce empty result",
			source:      `{}`,
			target:      `{}`,
			wantAdded:   nil,
			wantRemoved: nil,
			wantChanged: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := CompareJSON(
				strings.NewReader(tc.source),
				strings.NewReader(tc.target),
			)

			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.ElementsMatch(t, tc.wantAdded, result.Added)
			assert.ElementsMatch(t, tc.wantRemoved, result.Removed)
			assert.ElementsMatch(t, tc.wantChanged, result.Changed)
		})
	}
}

func TestComparisonResult_IsEmpty(t *testing.T) {
	tests := []struct {
		name   string
		result ComparisonResult
		want   bool
	}{
		{
			name:   "zero value is empty",
			result: ComparisonResult{},
			want:   true,
		},
		{
			name:   "only Added populated is not empty",
			result: ComparisonResult{Added: []string{"x"}},
			want:   false,
		},
		{
			name:   "only Removed populated is not empty",
			result: ComparisonResult{Removed: []string{"x"}},
			want:   false,
		},
		{
			name:   "only Changed populated is not empty",
			result: ComparisonResult{Changed: []string{"x"}},
			want:   false,
		},
		{
			name: "all fields populated is not empty",
			result: ComparisonResult{
				Added:   []string{"a"},
				Removed: []string{"b"},
				Changed: []string{"c"},
			},
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.result.IsEmpty())
		})
	}
}

func TestComparisonResult_Total(t *testing.T) {
	tests := []struct {
		name   string
		result ComparisonResult
		want   int
	}{
		{
			name:   "zero value totals zero",
			result: ComparisonResult{},
			want:   0,
		},
		{
			name:   "only Added counts correctly",
			result: ComparisonResult{Added: []string{"a", "b"}},
			want:   2,
		},
		{
			name:   "only Removed counts correctly",
			result: ComparisonResult{Removed: []string{"a"}},
			want:   1,
		},
		{
			name:   "only Changed counts correctly",
			result: ComparisonResult{Changed: []string{"a", "b", "c"}},
			want:   3,
		},
		{
			name: "all fields summed correctly",
			result: ComparisonResult{
				Added:   []string{"a", "b"},
				Removed: []string{"c"},
				Changed: []string{"d", "e"},
			},
			want: 5,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.result.Total())
		})
	}
}
