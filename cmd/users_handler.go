package cmd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func HandleUsers(state *state, cmd CliCommand) error {
	if len(cmd.Args) != 1 {
		return errors.New("Usage: gator users")
	}

	users, err := state.DatabaseQueries.GetUsers(context.Background())

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}

	for _, user := range users {
		prefix := "-"
		sufix := ""

		if user.Name == state.Config.CurrentUserName {
			sufix = " (current)"
		}

		fmt.Printf("%s %s%s\n", prefix, user.Name, sufix)
	}

	return nil
}
