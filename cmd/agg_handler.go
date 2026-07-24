package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/thalesraymond/gator/internal/database"
	"github.com/thalesraymond/gator/internal/feed"
)

func HandleAggregator(state *state, cmd CliCommand) error {
	for {
		feedToFetch, err := state.DatabaseQueries.GetNextFeedToFetch(context.Background())

		if err != nil {
			return err
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

			fmt.Printf("Item: %s\n", item.Title)
			fmt.Printf("Link: %s\n", item.Link)
			fmt.Printf("Description: %s\n", item.Description)

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
