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
	if err := c.CmdMap[cmd.Name](s, cmd); err != nil {
		return err
	}
	return nil
}

func (c *Commands) Register(name string, f func(*config.State, Command) error) {
	c.CmdMap[name] = f
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}
