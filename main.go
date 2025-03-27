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
	// commandRegistry := newCommands()

	cmds := os.Args[1:]
	if len(cmds) < 2 {
		fmt.Println("error not enough commands/arguments given")
		os.Exit(1)
	}

}
