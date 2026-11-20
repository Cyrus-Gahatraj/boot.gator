package main

import (
	"testing"
	"time"
)

func TestParsePublishedAt(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  time.Time
	}{
		{"RFC1123Z", "Sun, 08 Jan 2023 00:00:00 +0000", time.Date(2023, 1, 8, 0, 0, 0, 0, time.UTC)},
		{"RFC1123 with zone name", "Sun, 08 Jan 2023 00:00:00 GMT", time.Date(2023, 1, 8, 0, 0, 0, 0, time.UTC)},
		{"RFC3339", "2023-01-08T00:00:00Z", time.Date(2023, 1, 8, 0, 0, 0, 0, time.UTC)},
		{"single digit day", "Sun, 8 Jan 2023 00:00:00 -0700", time.Date(2023, 1, 8, 0, 0, 0, 0, time.FixedZone("", -7*3600))},
		{"date and time", "2023-01-08 00:00:00", time.Date(2023, 1, 8, 0, 0, 0, 0, time.UTC)},
		{"date only", "2023-01-08", time.Date(2023, 1, 8, 0, 0, 0, 0, time.UTC)},
		{"whitespace", "  2023-01-08T00:00:00Z  ", time.Date(2023, 1, 8, 0, 0, 0, 0, time.UTC)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePublishedAt(tt.input)
			if err != nil {
				t.Fatalf("parsePublishedAt(%q) returned error: %v", tt.input, err)
			}
			if !got.Equal(tt.want) {
				t.Errorf("parsePublishedAt(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParsePublishedAtInvalid(t *testing.T) {
	if _, err := parsePublishedAt("not a date"); err == nil {
		t.Fatal("expected an error for an unparseable date")
	}
}
