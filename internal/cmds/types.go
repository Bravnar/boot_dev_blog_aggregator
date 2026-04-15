// Package cmds holds the collection of commands to be used for gator
package cmds

import "github.com/Bravnar/gator/internal/config"

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	CmdMap map[string]func(*config.State, Command) error
}

func (c *Commands) Run(s *config.State, cmd Command) error {
	return nil
}

func (c *Commands) Register(name string, f func(*config.State, Command) error) {
}
