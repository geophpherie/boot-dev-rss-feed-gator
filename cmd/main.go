package main

import (
	"boot-dev-gator/internal/config"
	"boot-dev-gator/internal/database"
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type State struct {
	config *config.Config
	db     *database.Queries
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		panic(err)
	}

	dbQueries := database.New(db)

	state := State{
		config: &cfg,
		db:     dbQueries,
	}

	cmds := commands{
		commands: make(map[string]func(*State, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerAllFeeds)

	args := os.Args
	if len(args) < 2 {
		log.Fatal("Not enough arguments")
	}

	cmdName := args[1]
	cmdArgs := args[2:]

	err = cmds.run(&state, command{name: cmdName, args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	} else {
		os.Exit(0)
	}

}
