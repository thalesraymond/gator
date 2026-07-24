package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/thalesraymond/gator/internal/database"
	"github.com/thalesraymond/gator/internal/feed"
)

var supportedPubDateLayouts = []string{
	time.RFC1123Z,
	time.RFC1123,
	time.RFC822Z,
	time.RFC822,
	time.RFC3339,
	time.RFC3339Nano,
	"Mon, 02 Jan 2006 15:04:05 -0700",
	"Mon, 2 Jan 2006 15:04:05 -0700",
	"2006-01-02T15:04:05-0700",
	"2006-01-02 15:04:05 -0700",
}

func parsePublishedAt(raw string) (time.Time, error) {
	dateStr := strings.TrimSpace(raw)

	if dateStr == "" {
		return time.Time{}, fmt.Errorf("empty published date")
	}

	for _, layout := range supportedPubDateLayouts {
		parsed, err := time.Parse(layout, dateStr)

		if err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("unsupported published date format: %q", dateStr)
}

func isDuplicatePostURLError(err error) bool {
	pqErr, ok := err.(*pq.Error)

	if !ok {
		return false
	}

	return pqErr.Code == "23505" && pqErr.Constraint == "posts_url_key"
}

func HandleAggregator(state *state, cmd CliCommand) error {
	for {
		feedToFetch, err := state.DatabaseQueries.GetNextFeedToFetch(context.Background())

		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("No feeds to fetch, sleeping for 30 seconds...")
				time.Sleep(30 * time.Second)
				continue
			}

			return fmt.Errorf("failed to get next feed to fetch: %w", err)
		}
		fmt.Println("============================================================")
		fmt.Printf(" ========== Fetching feed: %s ========== \n", feedToFetch.Url)
		fmt.Println("============================================================")

		rssFeed, err := feed.FetchFeed(context.Background(), feedToFetch.Url)

		if err != nil {
			fmt.Printf("Failed to fetch feed: %v\n", err)
			continue
		}

		for _, item := range rssFeed.Channel.Item {
			publishedAt, err := parsePublishedAt(item.PubDate)

			if err != nil {
				log.Printf("failed to parse published date for post %q (%q): %v", item.Title, item.PubDate, err)
				continue
			}

			description := sql.NullString{}
			trimmedDescription := strings.TrimSpace(item.Description)

			if trimmedDescription != "" {
				description = sql.NullString{
					String: trimmedDescription,
					Valid:  true,
				}
			}

			_, err = state.DatabaseQueries.CreatePost(context.Background(), database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   sql.NullTime{},
				Title:       strings.TrimSpace(item.Title),
				Url:         strings.TrimSpace(item.Link),
				Description: description,
				PublishedAt: publishedAt,
				FeedID:      feedToFetch.ID,
			})

			if err != nil {
				if isDuplicatePostURLError(err) {
					continue
				}

				log.Printf("failed to save post %q (%s): %v", item.Title, item.Link, err)
			}

		}

		params := database.MarkFeedFetchedParams{
			ID: feedToFetch.ID,
			LastFetchedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
		}

		if err := state.DatabaseQueries.MarkFeedFetched(context.Background(), params); err != nil {
			fmt.Printf("Failed to mark feed as fetched: %v\n", err)
		}

		time.Sleep(10 * time.Second)

	}
}
