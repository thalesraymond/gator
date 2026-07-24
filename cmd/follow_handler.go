package cmd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/thalesraymond/gator/internal/database"
)

func HandleFollow(state *state, cmd CliCommand, authUser *database.User) error {
	if len(cmd.Args) != 2 {
		return errors.New("usage: gator follow <url>")
	}

	url := cmd.Args[1]

	feed, err := state.DatabaseQueries.GetFeedByURL(context.Background(), url)

	if err != nil {
		return err
	}

	followID := uuid.New()

	createdFollow, err := state.DatabaseQueries.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        followID,
		CreatedAt: time.Now(),
		UserID:    authUser.ID,
		FeedID:    feed.ID,
	})

	if err != nil {
		return err
	}

	fmt.Printf("Followed feed: %s\n", createdFollow.FeedID)
	fmt.Printf("Follow ID: %s\n", createdFollow.ID)
	fmt.Printf("Follow created at: %s\n", createdFollow.CreatedAt)
	fmt.Printf("Follow user name: %s\n", createdFollow.UserName)
	fmt.Printf("Follow feed name: %s\n", createdFollow.FeedName)

	return nil
}
