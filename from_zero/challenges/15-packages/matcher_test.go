// This file tests the matcher package directly (unit tests).
//
// IMPORTANT: Copy this file to internal/matcher/matcher_test.go
// in your project. It will NOT work from this location — it must
// live alongside matcher.go in the matcher package.
//
// Run: go test -v ./internal/matcher/

package matcher_test

import (
	"search/internal/matcher"
	"testing"
)

// ---------- Match function ----------

func TestMatch(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		line    string
		want    bool
	}{
		{
			name:    "exact match",
			pattern: "Judas Priest",
			line:    "Judas Priest",
			want:    true,
		},
		{
			name:    "substring match at start",
			pattern: "Judas",
			line:    "Judas Priest",
			want:    true,
		},
		{
			name:    "substring match at end",
			pattern: "Priest",
			line:    "Judas Priest",
			want:    true,
		},
		{
			name:    "substring match in middle",
			pattern: "das P",
			line:    "Judas Priest",
			want:    true,
		},
		{
			name:    "no match",
			pattern: "Metallica",
			line:    "Judas Priest",
			want:    false,
		},
		{
			name:    "case sensitive — lowercase",
			pattern: "priest",
			line:    "Judas Priest",
			want:    false,
		},
		{
			name:    "case sensitive — uppercase",
			pattern: "PRIEST",
			line:    "Judas Priest",
			want:    false,
		},
		{
			name:    "empty pattern matches everything",
			pattern: "",
			line:    "Judas Priest",
			want:    true,
		},
		{
			name:    "empty line — no match",
			pattern: "test",
			line:    "",
			want:    false,
		},
		{
			name:    "both empty",
			pattern: "",
			line:    "",
			want:    true,
		},
		{
			name:    "special characters — slash",
			pattern: "/",
			line:    "AC/DC",
			want:    true,
		},
		{
			name:    "special characters — apostrophe",
			pattern: "'",
			line:    "Guns N' Roses",
			want:    true,
		},
		{
			name:    "special characters — ampersand",
			pattern: "&",
			line:    "Y&T",
			want:    true,
		},
		{
			name:    "special characters — dot",
			pattern: ".",
			line:    "W.A.S.P.",
			want:    true,
		},
		{
			name:    "single character match",
			pattern: "K",
			line:    "Kiss",
			want:    true,
		},
		{
			name:    "pattern longer than line",
			pattern: "Judas Priest is great",
			line:    "Judas Priest",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matcher.Match(tt.pattern, tt.line)
			if got != tt.want {
				t.Errorf("Match(%q, %q) = %v, want %v", tt.pattern, tt.line, got, tt.want)
			}
		})
	}
}

// ---------- SearchLines function ----------

func TestSearchLines(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		lines   []string
		want    []string
	}{
		{
			name:    "single match",
			pattern: "Priest",
			lines:   []string{"Judas Priest", "AC/DC", "Black Sabbath"},
			want:    []string{"Judas Priest"},
		},
		{
			name:    "multiple matches",
			pattern: "Bad",
			lines:   []string{"Bad English", "Foreigner", "Bad Company", "Boston"},
			want:    []string{"Bad English", "Bad Company"},
		},
		{
			name:    "no matches",
			pattern: "zzz",
			lines:   []string{"Judas Priest", "AC/DC", "Black Sabbath"},
			want:    []string{},
		},
		{
			name:    "all match",
			pattern: "a",
			lines:   []string{"banana", "avocado", "apricot"},
			want:    []string{"banana", "avocado", "apricot"},
		},
		{
			name:    "empty lines slice",
			pattern: "test",
			lines:   []string{},
			want:    []string{},
		},
		{
			name:    "preserves order",
			pattern: "White",
			lines:   []string{"Whitesnake", "Great White", "White Lion", "Whitecross"},
			want:    []string{"Whitesnake", "Great White", "White Lion", "Whitecross"},
		},
		{
			name:    "case sensitive",
			pattern: "bad",
			lines:   []string{"Bad English", "Bad Company"},
			want:    []string{},
		},
		{
			name:    "partial matches",
			pattern: "ap",
			lines:   []string{"apple", "banana", "cherry", "apricot", "blueberry", "avocado"},
			want:    []string{"apple", "apricot"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matcher.SearchLines(tt.pattern, tt.lines)

			// Handle nil vs empty slice — both are acceptable for "no results"
			if len(got) == 0 && len(tt.want) == 0 {
				return
			}

			if len(got) != len(tt.want) {
				t.Fatalf("SearchLines(%q, ...) returned %d results, want %d\n  got:  %v\n  want: %v",
					tt.pattern, len(got), len(tt.want), got, tt.want)
			}

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("SearchLines(%q, ...)[%d] = %q, want %q",
						tt.pattern, i, got[i], tt.want[i])
				}
			}
		})
	}
}

// ---------- SearchLines with nil input ----------

func TestSearchLinesNilInput(t *testing.T) {
	got := matcher.SearchLines("test", nil)
	if len(got) != 0 {
		t.Errorf("SearchLines with nil input should return empty, got: %v", got)
	}
}
