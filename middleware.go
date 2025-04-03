package main

import (
	"context"

	"github.com/ChipsAhoyEnjoyer/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, c command) error {
		usr, err := s.db.GetUser(
			context.Background(),
			s.cfg.CurrentUsername,
		)
		if err != nil {
			return err
		}
		err = handler(s, c, usr)
		if err != nil {
			return err
		}
		return nil
	}
}
