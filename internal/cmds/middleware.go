package cmds

import (
	"context"

	"github.com/Bravnar/gator/internal/config"
	"github.com/Bravnar/gator/internal/database"
)

func MiddlewareLoggedIn(handler func(s *config.State, cmd Command, user database.User) error) func(*config.State, Command) error {
	return func(s *config.State, cmd Command) error {
		user, err := s.DB.GetUser(context.Background(), s.Conf.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}
