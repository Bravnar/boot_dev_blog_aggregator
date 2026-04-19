package cmds

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
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

func HandlerUsers(s *config.State, cmd Command) error {
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

func HandlerAddFeed(s *config.State, cmd Command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage ./gator \"blog title\" \"blog url\"")
	}

	currentUser := s.Conf.CurrentUserName
	userUUID, err := s.DB.GateUserUUID(context.Background(), currentUser)
	if err != nil {
		return err
	}

	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    userUUID,
	}
	feed, err := s.DB.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}
	fmt.Printf("ID: %v\n", feed.ID)
	fmt.Printf("CreatedAt: %v\n", feed.CreatedAt)
	fmt.Printf("UpdatedAt: %v\n", feed.UpdatedAt)
	fmt.Printf("Name: %v\n", feed.Name)
	fmt.Printf("Url: %v\n", feed.Url)
	fmt.Printf("UserID: %v\n", feed.UserID)
	return nil
}

func prettyPrintXML(x *RSSFeed) {
	fmt.Printf("Title: %s\nDescription: %s\nLink: %s", x.Channel.Title, x.Channel.Description, x.Channel.Link)
	fmt.Println("Items:")
	for _, item := range x.Channel.Item {
		fmt.Printf(" * Title: %s\n", item.Title)
		fmt.Printf(" * Description: %s\n", item.Description)
	}
}

func HandlerAgg(s *config.State, cmd Command) error {
	xmlFeed, err := FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	prettyPrintXML(xmlFeed)
	return nil
}

// -------------------  commands to fetch the XML Feed

func cleanUpFeed(feed *RSSFeed) {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i, item := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(item.Title)
		feed.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}
}

func decodeXML(toDecode io.Reader) (*RSSFeed, error) {
	var feed *RSSFeed
	decoder := xml.NewDecoder(toDecode)
	if err := decoder.Decode(&feed); err != nil {
		return feed, err
	}
	cleanUpFeed(feed)
	return feed, nil
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	var feed *RSSFeed
	client := &http.Client{}

	req, err := http.NewRequestWithContext(context.Background(), "GET", feedURL, nil)
	if err != nil {
		return feed, err
	}
	req.Header.Set("User-Agent", "gator")

	resp, err := client.Do(req)
	if err != nil {
		return feed, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return feed, fmt.Errorf("bad status code: %v", resp.StatusCode)
	}

	byteData, err := io.ReadAll(resp.Body)
	if err != nil {
		return feed, err
	}

	return decodeXML(bytes.NewReader(byteData))
}
