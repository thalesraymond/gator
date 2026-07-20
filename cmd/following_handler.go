package cmd

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func HandleFollowing(state *state, cmd CliCommand) error {
	if len(cmd.Args) != 1 {
		return errors.New("usage: gator following")
	}

	authUser, err := state.DatabaseQueries.GetUserByName(context.Background(), state.Config.CurrentUserName)

	if err != nil {
		return err
	}

	following, err := state.DatabaseQueries.GetFeedFollowsForUser(context.Background(), authUser.ID)

	if err != nil {
		return err
	}

	if len(following) == 0 {
		fmt.Println("You are not following any feeds.")
		return nil
	}

	fmt.Println("You are following these feeds:")
	for _, follow := range following {
		fmt.Printf("- %s (followed on %s)\n", follow.FeedName, follow.CreatedAt.Format(time.RFC1123))
	}

	return nil
}
