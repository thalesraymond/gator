package cmd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/thalesraymond/gator/internal/database"
)

func HandleAddFeed(state *state, cmd CliCommand, authUser *database.User) error {
	if len(cmd.Args) != 3 {
		return errors.New("usage: gator add_feed <name> <url>")
	}

	name := cmd.Args[1]
	url := cmd.Args[2]

	feedID := uuid.New()

	createdFeed, err := state.DatabaseQueries.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        feedID,
		Name:      name,
		Url:       url,
		UserID:    authUser.ID,
		CreatedAt: time.Now(),
	})

	if err != nil {
		return err
	}

	_, err = state.DatabaseQueries.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UserID:    authUser.ID,
		FeedID:    createdFeed.ID,
	})

	if err != nil {
		return err
	}

	fmt.Printf("Feed created: %s (%s)\n", createdFeed.Name, createdFeed.Url)

	return nil
}
