package cmds

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Bravnar/gator/internal/config"
	"github.com/Bravnar/gator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
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

	for _, item := range rssFeed.Channel.Item {
		pubDate := item.PubDate
		t, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", pubDate)
		if err != nil {
			t, err = time.Parse("2006-01-02T15:04:05Z", pubDate)
			if err != nil {
				return err
			}
		}
		postParams := database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: t,
			FeedID:      nextFeedToFetch.ID,
		}
		post, err := s.DB.CreatePost(context.Background(), postParams)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				continue
			}
			log.Printf("error creating post: %v", err)
			continue
		}
		fmt.Printf("Post: %s\nSuccessfully saved to DB\n", post.Title)
	}

	return nil
}
