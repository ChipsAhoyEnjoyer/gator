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
	c.register("register", handlerRegister)
	c.register("reset", handlerReset)
	c.register("users", handlerGetUsers)
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

func handlerGetUsers(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("error: too many arguments given; users expects zero arguments")
	}
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error: could not retreive users from db \n%v", err)
	}
	for _, user := range users {
		if user.String == s.cfg.CurrentUsername {
			fmt.Printf("* %v (current)\n", user.String)
		} else {
			fmt.Printf("* %v\n", user.String)
		}
	}
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) >= 2 {
		return fmt.Errorf("error: too many arguments given; login expects one(username) argument")
	}
	username := cmd.args[0]
	if !userExists(s, username) {
		return fmt.Errorf("error: user not registered")
	}
	err := s.cfg.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Printf("Now logged in as %v\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) >= 2 {
		return fmt.Errorf("error: too many arguments given; register expects one(username) argument")
	}
	username := cmd.args[0]
	if userExists(s, username) {
		return fmt.Errorf("error: user exists")
	}
	u, err := s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name: sql.NullString{
				String: username,
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

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("error: too many arguments given; reset expects zero arguments")
	}
	usersDeleted, err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error: users table reset unsuccessful \n%v", err)
	}
	fmt.Printf("Deleted %v user(s)\n", usersDeleted)
	return nil
}

func userExists(s *state, name string) bool {
	_, err := s.db.GetUser(
		context.Background(),
		sql.NullString{
			String: name,
			Valid:  true,
		},
	)
	return err == nil
}
