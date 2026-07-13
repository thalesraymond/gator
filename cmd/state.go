package cmd

import (
	"database/sql"

	"github.com/thalesraymond/gator/internal/config"
	"github.com/thalesraymond/gator/internal/database"
)

type state struct {
	Config          *config.Config
	DatabaseQueries *database.Queries
}

func NewState() (*state, error) {
	config, err := config.ReadConfig()

	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", config.DbUrl)
	if err != nil {
		return nil, err
	}

	return &state{Config: config, DatabaseQueries: database.New(db)}, nil
}
