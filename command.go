package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ChipsAhoyEnjoyer/gator/internal/database"
	"github.com/google/uuid"
)

type command struct {
	name string
	args []string
}

type commands struct {
	registry map[string]func(*state, command) error
}

func newCommands() *commands {
	c := commands{
		registry: make(map[string]func(*state, command) error),
	}
	c.register("login", handlerLogin)
	return &c
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registry[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	if function, ok := c.registry[cmd.name]; !ok {
		return fmt.Errorf("error: command '%v' does not exist", cmd.name)
	} else {
		err := function(s, cmd)
		if err != nil {
			return err
		}
	}
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) > 2 {
		return fmt.Errorf("error: too many arguments given; login expects one(username) argument")
	}
	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Now logged in as %v\n", cmd.args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) > 2 {
		return fmt.Errorf("error: too many arguments given; register expects one(username) argument")
	}
	u, err := s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name: sql.NullString{
				String: cmd.args[1],
				Valid:  true,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("error: could not register user to database\n%v", err)
	}
	err = s.cfg.SetUser(u.Name.String)
	if err != nil {
		return fmt.Errorf("error: user registered but not logged in\n%v", err)
	}
	fmt.Printf(
		"User '%v' created \nCreated: %v \nUpdated: %v \nid: %v \n",
		u.Name.String,
		u.CreatedAt,
		u.UpdatedAt,
		u.ID,
	)
	return nil
}
