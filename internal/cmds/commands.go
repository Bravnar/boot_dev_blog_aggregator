package cmds

import (
	"context"
	"fmt"
	"time"

	"github.com/Bravnar/gator/internal/config"
	"github.com/Bravnar/gator/internal/database"
	"github.com/google/uuid"
)

func HandlerLogic(s *config.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("expected user login as argument")
	}
	login := cmd.Args[0]
	_, err := s.DB.GetUser(context.Background(), login)
	if err != nil {
		return fmt.Errorf("user does not exist")
	}
	if err := s.Conf.SetUser(login); err != nil {
		return err
	}
	fmt.Printf("Username: %s - has been set.\n", login)
	return nil
}

func HandlerRegister(s *config.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("expected a name to be added as argument")
	}
	name := cmd.Args[0]
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}
	_, err := s.DB.CreateUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("user already exists")
	}
	fmt.Printf("User %s successfully created\n", name)
	if err := s.Conf.SetUser(name); err != nil {
		return err
	}
	return nil
}

func HandlerReset(s *config.State, cmd Command) error {
	if err := s.DB.DeleteAllUsers(context.Background()); err != nil {
		return fmt.Errorf("failed to delete all users")
	}
	fmt.Println("deleted all users from database")
	return nil
}

func HandleUsers(s *config.State, cmd Command) error {
	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to fetch users from database")
	}
	current := s.Conf.CurrentUserName
	for _, user := range users {
		toPrint := "* " + user.Name
		if current == user.Name {
			toPrint += " (current)"
		}
		fmt.Println(toPrint)
	}
	return nil
}
