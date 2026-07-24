package cmd

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/thalesraymond/gator/internal/database"
)

const defaultBrowseLimit int32 = 2

func HandleBrowse(state *state, cmd CliCommand, authUser *database.User) error {
	if len(cmd.Args) > 2 {
		return errors.New("usage: gator browse [limit]")
	}

	limit := defaultBrowseLimit

	if len(cmd.Args) == 2 {
		parsedLimit, err := strconv.Atoi(cmd.Args[1])

		if err != nil || parsedLimit < 1 {
			return errors.New("limit must be a positive integer")
		}

		limit = int32(parsedLimit)
	}

	posts, err := state.DatabaseQueries.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: authUser.ID,
		Limit:  limit,
	})

	if err != nil {
		return err
	}

	if len(posts) == 0 {
		fmt.Println("No posts found.")
		return nil
	}

	for _, post := range posts {
		fmt.Println("------------------------------")
		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("URL: %s\n", post.Url)
		fmt.Printf("Published: %s\n", post.PublishedAt.Format("2006-01-02 15:04:05 MST"))

		if post.Description.Valid {
			fmt.Printf("Description: %s\n", post.Description.String)
		}
	}

	return nil
}
