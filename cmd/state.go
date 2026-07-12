package cmd

import (
	"github.com/thalesraymond/gator/internal/config"
)

type state struct {
	Config *config.Config
}

func NewState() (*state, error) {
	config, err := config.ReadConfig()
	if err != nil {
		return nil, err
	}

	return &state{Config: config}, nil
}
