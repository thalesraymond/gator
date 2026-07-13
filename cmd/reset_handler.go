package cmd

import (
	"context"
	"database/sql"
	"errors"
)

func HandleReset(state *state, cmd CliCommand) error {
	if len(cmd.Args) != 1 {
		return errors.New("Usage: gator reset")
	}

	err := state.DatabaseQueries.DeleteAllUsers(context.Background())

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}

	return nil
}
