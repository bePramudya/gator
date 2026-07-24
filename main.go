package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/bePramudya/gator/internal/config"
	"github.com/bePramudya/gator/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	dbQueries := database.New(db)

	appState := state{
		db:  dbQueries,
		cfg: &cfg,
	}

	var appCommands = commands{
		handlers: map[string]func(*state, command) error{},
	}

	appCommands.register("login", handlerLogin)
	appCommands.register("register", handlerRegister)
	appCommands.register("reset", handlerReset)
	appCommands.register("users", handlerGetUsers)
	appCommands.register("agg", handlerAggregate)
	appCommands.register("addfeed", handlerAddFeed)
	appCommands.register("feeds", handlerListFeeds)

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
