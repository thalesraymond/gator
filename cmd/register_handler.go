package cmd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/thalesraymond/gator/internal/database"
)

func HandleRegister(state *state, cmd CliCommand) error {
	if len(cmd.Args) != 2 {
		return errors.New("usage: gator register <username>")
	}

	username := cmd.Args[1]

	user, err := state.DatabaseQueries.GetUserByName(context.Background(), username)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}

	if user.Name == username {
		return errors.New("user already exists")
	}

	now := time.Now()

	createdUser, err := state.DatabaseQueries.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      username,
		CreatedAt: now,
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		return err
	}

	if err := state.Config.SetUser(createdUser.Name); err != nil {
		return err
	}

	fmt.Println("User created successfully")
	fmt.Printf("Logged in as: %s\n", createdUser.Name)

	return nil
}
