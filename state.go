package main

import (
	"fmt"

	"github.com/ChipsAhoyEnjoyer/gator/internal/config"
)

type state struct {
	config *config.Config
}

func createStateInstance() *state {
	return &state{}
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("error login expects a username as an argument")
	}
	err := s.config.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Now logged in as %v", cmd.args[0])
	return nil
}
