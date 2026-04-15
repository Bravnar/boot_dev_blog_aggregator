package cmds

import (
	"fmt"

	"github.com/Bravnar/gator/internal/config"
)

func HandlerLogic(s *config.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("expected user login as argument")
	}
	login := cmd.Args[0]
	if err := s.Conf.SetUser(login); err != nil {
		return err
	}
	fmt.Printf("Username: %s - has been set.", login)
	return nil
}
