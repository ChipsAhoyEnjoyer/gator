package main

import (
	"github.com/ChipsAhoyEnjoyer/gator/internal/config"
	"github.com/ChipsAhoyEnjoyer/gator/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func createStateInstance() *state {
	return &state{
		cfg: &config.Config{},
	}
}
