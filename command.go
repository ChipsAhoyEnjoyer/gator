package main

import "fmt"

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
