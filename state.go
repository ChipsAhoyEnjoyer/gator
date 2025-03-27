package main

import (
	"github.com/ChipsAhoyEnjoyer/gator/internal/config"
)

type state struct {
	config *config.Config
}

func createStateInstance() *state {
	return &state{
		config: &config.Config{},
	}
}
