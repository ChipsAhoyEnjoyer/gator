package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/ChipsAhoyEnjoyer/gator/internal/config"
	"github.com/ChipsAhoyEnjoyer/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	gatorState := createStateInstance()
	newConfig, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	gatorState.cfg = newConfig
	commandRegistry := newCommands()
	db, err := sql.Open("postgres", gatorState.cfg.DBUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	gatorState.db = database.New(db)
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

// TODO: fix this clean input to accept commands with no args and 1+ arg(s)
func cleanInput(input []string) (*command, error) {
	if len(input) == 2 {
		if input[1] == "reset" ||
			input[1] == "users" ||
			input[1] == "agg" ||
			input[1] == "following" ||
			input[1] == "feeds" {
			return &command{name: input[1]}, nil
		}
	}
	if len(input) < 3 {
		return nil, fmt.Errorf("error: not enough commands/arguments given")
	}
	cmd := command{
		name: input[1],
		args: input[2:],
	}
	return &cmd, nil
}
