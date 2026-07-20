package cmd

import (
	"context"
	"errors"
	"fmt"
)

func HandleLogin(state *state, cmd CliCommand) error {
	if len(cmd.Args) != 2 {
		return errors.New("usage: gator login <username>")
	}

	username := cmd.Args[1]

	user, err := state.DatabaseQueries.GetUserByName(context.Background(), username)

	if err != nil {
		return err
	}

	if err := state.Config.SetUser(user.Name); err != nil {
		return err
	}

	fmt.Printf("Logged in as: %s\n", username)

	return nil
}
