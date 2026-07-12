package cmd

import (
	"errors"
	"fmt"
)

func HandleLogin(state *state, cmd CliCommand) error {
	if len(cmd.Args) != 2 {
		return errors.New("Usage: gator login <username>")
	}

	username := cmd.Args[1]

	if err := state.Config.SetUser(username); err != nil {
		return err
	}

	fmt.Printf("Logged in as: %s\n", username)

	return nil
}
