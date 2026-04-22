package cmds

import (
	"context"
	"fmt"
	"time"

	"github.com/Bravnar/gator/internal/config"
	"github.com/Bravnar/gator/internal/database"
)

func scrapeFeeds(s *config.State) error {
	nextFeedToFetch, err := s.DB.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	params := database.MarkFeedFetchedParams{
		UpdatedAt: time.Now(),
		ID:        nextFeedToFetch.ID,
	}
	if err = s.DB.MarkFeedFetched(context.Background(), params); err != nil {
		return err
	}

	rssFeed, err := FetchFeed(context.Background(), nextFeedToFetch.Url)
	if err != nil {
		return err
	}
	fmt.Printf("\n### ------------- | NEW FETCH | ----------------- ###\n")
	fmt.Printf("Items of %s\n", rssFeed.Channel.Title)
	for i, item := range rssFeed.Channel.Item {
		fmt.Printf("\t%d. %s\n", i+1, item.Title)
	}

	return nil
}
