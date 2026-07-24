package cmd

import (
	"context"
	"fmt"

	"github.com/thalesraymond/gator/internal/database"
)

func HandlerListFeeds(state *state, cmd CliCommand, authUser *database.User) error {
	feeds, err := state.DatabaseQueries.ListAllFeeds(context.Background())

	if err != nil {
		return err
	}

	for _, f := range feeds {
		fmt.Printf("Feed: %s (%s) for user %s\n", f.Name, f.Url, f.UserName)
	}

	return nil
}
