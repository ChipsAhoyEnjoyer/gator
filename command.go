package main

import (
	"fmt"
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
	c.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	c.register("following", middlewareLoggedIn(handlerFollowing))
	c.register("follow", middlewareLoggedIn(handlerFollow))
	c.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	c.register("agg", middlewareLoggedIn(handlerAgg))
	c.register("browse", middlewareLoggedIn(handlerBrowse))
	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	c.register("reset", handlerReset)
	c.register("users", handlerUsers)
	c.register("feeds", handlerFeeds)
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
