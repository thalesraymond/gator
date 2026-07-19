package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/thalesraymond/gator/internal/feed"
)

func HandleAggregator(state *state, cmd CliCommand) error {
	// if len(cmd.Args) != 2 {
	// 	return errors.New("Usage: gator agg <url>")
	// }

	url := "https://www.wagslane.dev/index.xml"

	feed, err := feed.FetchFeed(context.Background(), url)
	if err != nil {
		return err
	}

	//print the whole struct as json
	fmt.Println("Feed fetched successfully")
	feedJson, err := json.MarshalIndent(feed, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(feedJson))

	return nil
}
