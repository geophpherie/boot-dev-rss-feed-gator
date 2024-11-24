package main

import (
	"boot-dev-gator/internal/database"
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}

	req.Header.Add("User-Agent", "gator")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, err
	}

	rssFeed := RSSFeed{}
	if err = xml.Unmarshal(body, &rssFeed); err != nil {
		return &RSSFeed{}, err
	}

	return &rssFeed, nil
}

func scrapeFeeds(s *State) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	now := time.Now()
	args := database.MarkFetchedFeedParams{
		LastFetchedAt: sql.NullTime{Time: now, Valid: true},
		UpdatedAt:     now,
		ID:            nextFeed.ID,
	}
	s.db.MarkFetchedFeed(context.Background(), args)

	feed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return err
	}

	cleanFeed := RSSFeed{}

	cleanFeed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	cleanFeed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	cleanFeed.Channel.Link = feed.Channel.Link

	for _, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)

		cleanFeed.Channel.Item = append(cleanFeed.Channel.Item, item)
		fmt.Println("creating post")
		t, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", item.PubDate)
		if err != nil {
			fmt.Print("Unable to parse time", item.PubDate)
			continue
		}
		now := time.Now()
		args := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   now,
			UpdatedAt:   now,
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: sql.NullTime{Time: t, Valid: true},
			FeedID:      nextFeed.ID,
		}
		_, err = s.db.CreatePost(context.Background(), args)
		if err != nil {
			fmt.Println("Unable to create Post")
			fmt.Println(err)
			return err
		}
	}

	return nil

}
