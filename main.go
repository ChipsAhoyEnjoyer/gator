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

func cleanInput(input []string) (*command, error) {
	if len(input) < 2 {
		return nil, fmt.Errorf("usage: cli <command> [args...]")
	}
	var cmd command
	cmd.name = os.Args[1]
	cmd.args = os.Args[2:]
	return &cmd, nil
}
