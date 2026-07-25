package feed

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchFeedTable(t *testing.T) {
	tests := []struct {
		name         string
		feedURL      string
		responseBody string
		expectErr    bool
		wantTitle    string
		wantDesc     string
		wantItemTitle string
		checkUA      bool
	}{
		{
			name: "fetches and unescapes feed content",
			responseBody: `<?xml version="1.0" encoding="UTF-8"?>
<rss>
	<channel>
		<title>Go &amp; Rust</title>
		<link>https://example.com</link>
		<description>Dev &lt;News&gt;</description>
		<item>
			<title>A &amp; B</title>
			<link>https://example.com/posts/1</link>
			<description>X &amp; Y</description>
			<pubDate>Mon, 02 Jan 2006 15:04:05 -0700</pubDate>
		</item>
	</channel>
</rss>`,
			wantTitle:     "Go & Rust",
			wantDesc:      "Dev <News>",
			wantItemTitle: "A & B",
			checkUA:       true,
		},
		{
			name:         "returns error on invalid xml",
			responseBody: `<rss><channel><title>broken</title>`,
			expectErr:    true,
		},
		{
			name:      "returns error on invalid url",
			feedURL:   "://bad-url",
			expectErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			feedURL := tc.feedURL
			uaChecked := false

			if feedURL == "" {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if tc.checkUA {
						uaChecked = true
						if got := r.Header.Get("User-Agent"); got != "Gator RSS Reader" {
							t.Errorf("expected User-Agent %q, got %q", "Gator RSS Reader", got)
						}
					}

					w.Header().Set("Content-Type", "application/xml")
					_, _ = w.Write([]byte(tc.responseBody))
				}))
				defer server.Close()

				feedURL = server.URL
			}

			result, err := FetchFeed(context.Background(), feedURL)

			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}

				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if tc.checkUA && !uaChecked {
				t.Fatalf("expected User-Agent assertion to run")
			}

			if result.Channel.Title != tc.wantTitle {
				t.Fatalf("expected title %q, got %q", tc.wantTitle, result.Channel.Title)
			}

			if result.Channel.Description != tc.wantDesc {
				t.Fatalf("expected description %q, got %q", tc.wantDesc, result.Channel.Description)
			}

			if len(result.Channel.Item) != 1 {
				t.Fatalf("expected 1 item, got %d", len(result.Channel.Item))
			}

			if result.Channel.Item[0].Title != tc.wantItemTitle {
				t.Fatalf("expected item title %q, got %q", tc.wantItemTitle, result.Channel.Item[0].Title)
			}
		})
	}
}
