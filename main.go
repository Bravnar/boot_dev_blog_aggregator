package main

import (
	"fmt"
	"os"

	"github.com/Bravnar/gator/internal/cmds"
	"github.com/Bravnar/gator/internal/config"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("not enough arguments provided")
		os.Exit(1)
	}
	cmd := cmds.Command{
		Name: args[1],
		Args: args[2:],
	}
	configFile, err := config.ReadConfig()
	if err != nil {
		fmt.Printf("failed to read configuration file: %s", err)
	}
	state := config.State{
		Conf: &configFile,
	}
	commandsMap := cmds.Commands{
		CmdMap: make(map[string]func(*config.State, cmds.Command) error),
	}
	commandsMap.Register("login", cmds.HandlerLogic)
	if err := commandsMap.Run(&state, cmd); err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("state: %v\n", state.Conf.CurrentUserName)
}
