package cmd

import (
	"errors"
	"testing"
	"time"

	"github.com/lib/pq"
)

func TestParsePublishedAtTable(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantErr bool
	}{
		{
			name: "parses RFC1123Z",
			raw:  "Mon, 02 Jan 2006 15:04:05 -0700",
		},
		{
			name: "parses RFC3339",
			raw:  "2006-01-02T15:04:05Z",
		},
		{
			name: "parses and trims surrounding whitespace",
			raw:  "  Mon, 2 Jan 2006 15:04:05 -0700  ",
		},
		{
			name:    "rejects empty date",
			raw:     "    ",
			wantErr: true,
		},
		{
			name:    "rejects unknown format",
			raw:     "not-a-date",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := parsePublishedAt(tc.raw)

			if tc.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}

			if !tc.wantErr {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				if got.IsZero() {
					t.Fatalf("expected non-zero time for %q", tc.raw)
				}

				if got.After(time.Now().Add(1000 * time.Hour)) {
					t.Fatalf("unexpected parsed time far in the future: %v", got)
				}
			}
		})
	}
}

func TestIsDuplicatePostURLErrorTable(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil error is not duplicate",
			err:  nil,
			want: false,
		},
		{
			name: "non pq error is not duplicate",
			err:  errors.New("boom"),
			want: false,
		},
		{
			name: "pq unique violation on posts_url_key is duplicate",
			err: &pq.Error{
				Code:       "23505",
				Constraint: "posts_url_key",
			},
			want: true,
		},
		{
			name: "pq unique violation on other constraint is not duplicate",
			err: &pq.Error{
				Code:       "23505",
				Constraint: "other_constraint",
			},
			want: false,
		},
		{
			name: "pq other code on posts_url_key is not duplicate",
			err: &pq.Error{
				Code:       "22001",
				Constraint: "posts_url_key",
			},
			want: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := isDuplicatePostURLError(tc.err)

			if got != tc.want {
				t.Fatalf("expected %v, got %v", tc.want, got)
			}
		})
	}
}
