package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Bravnar/gator/internal/cmds"
	"github.com/Bravnar/gator/internal/config"
	"github.com/Bravnar/gator/internal/database"
	_ "github.com/lib/pq"
)

func fetchArgs() []string { // without executable
	args := os.Args
	if len(args) < 2 {
		fmt.Println("not enough arguments provided")
		os.Exit(1)
	}
	return args[1:]
}

func main() {
	cmd := cmds.Command{
		Name: fetchArgs()[0],
		Args: fetchArgs()[1:],
	}

	configFile, err := config.ReadConfig()
	if err != nil {
		fmt.Printf("failed to read configuration file: %s", err)
	}

	state := config.State{
		Conf: &configFile,
	}

	db, err := sql.Open("postgres", state.Conf.DBURL)
	if err != nil {
		fmt.Printf("failed to connect to db")
		return
	}
	dbQueries := database.New(db)
	state.DB = dbQueries

	commandsMap := cmds.Commands{
		CmdMap: make(map[string]func(*config.State, cmds.Command) error),
	}

	commandsMap.Register("login", cmds.HandlerLogic)
	commandsMap.Register("register", cmds.HandlerRegister)
	commandsMap.Register("reset", cmds.HandlerReset)
	commandsMap.Register("users", cmds.HandleUsers)
	if err := commandsMap.Run(&state, cmd); err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}
