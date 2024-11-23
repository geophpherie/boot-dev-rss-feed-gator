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

		fmt.Println(item.Title)
	}

	return nil

}
