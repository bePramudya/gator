package main

import (
	"fmt"
	"os"

	"github.com/bePramudya/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	appState := state{
		Config: &cfg,
	}

	var appCommands = commands{
		handlers: map[string]func(*state, command) error{},
	}

	appCommands.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("no command provided")
		os.Exit(1)
	}

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	err = appCommands.run(&appState, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
