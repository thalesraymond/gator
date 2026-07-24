package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/thalesraymond/gator/internal/database"
)

func HandleUnfollow(state *state, cmd CliCommand, authUser *database.User) error {
	if len(cmd.Args) != 2 {
		return errors.New("usage: gator unfollow <url>")
	}

	url := cmd.Args[1]

	feed, err := state.DatabaseQueries.GetFeedByURL(context.Background(), url)

	if err != nil {
		return err
	}

	err = state.DatabaseQueries.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: authUser.ID,
		FeedID: feed.ID,
	})

	if err != nil {
		return err
	}

	fmt.Printf("Unfollowed feed: %s\n", feed.ID)

	return nil
}
