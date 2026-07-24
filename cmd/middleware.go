package cmd

import (
	"context"

	"github.com/thalesraymond/gator/internal/database"
)

func MiddlewareLogin(next func(state *state, cmd CliCommand, authUser *database.User) error) func(state *state, cmd CliCommand) error {
	return func(state *state, cmd CliCommand) error {
		authUser, err := state.DatabaseQueries.GetUserByName(context.Background(), state.Config.CurrentUserName)

		if err != nil {
			return err
		}

		return next(state, cmd, &authUser)
	}
}
