package main

import (
	"fmt"
	"os"

	"github.com/ChipsAhoyEnjoyer/gator/internal/config"
)

func main() {
	gatorState := createStateInstance()
	newConfig, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	gatorState.config = newConfig
	commandRegistry := newCommands()
	cmd, err := cleanInput(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = commandRegistry.run(gatorState, *cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func cleanInput(input []string) (*command, error) {
	if len(input) < 3 {
		return nil, fmt.Errorf("error: not enough commands/arguments given")
	}
	cmd := command{
		name: input[1],
		args: input[2:],
	}
	return &cmd, nil
}
