package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/Yah-oo5726/blog-aggregator/internal/database"
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

func unescape(s *string) {
	*s = html.UnescapeString(*s)
}

func (feed *RSSFeed) unescapeText() {
	unescape(&feed.Channel.Title)
	unescape(&feed.Channel.Description)
	for i := range feed.Channel.Item {
		unescape(&feed.Channel.Item[i].Title)
		unescape(&feed.Channel.Item[i].Description)
	}
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "gator")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	feed := RSSFeed{}
	feedPtr := &feed
	err = xml.Unmarshal(data, feedPtr)
	if err != nil {
		return nil, err
	}
	feedPtr.unescapeText()
	return feedPtr, nil
}

func scrapeFeed(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	s.db.MarkFeedFetched(context.Background(), nextFeed.ID)
	feed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return err
	}
	for _, item := range feed.Channel.Item {
		publishedAt, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", item.PubDate)
		if err != nil {
			return err
		}
		err = s.db.CreatePost(context.Background(), database.CreatePostParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Title: item.Title, Url: item.Link, Description: item.Description, PublishedAt: publishedAt, FeedID: nextFeed.ID})
		if err != nil {
			return err
		}
	}
	return nil
}
