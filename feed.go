package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
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
		fmt.Println(err.Error())
		os.Exit(1)
	}
	request.Header.Set("User-Agent", "gator")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	feed := RSSFeed{}
	feedPtr := &feed
	err = xml.Unmarshal(data, feedPtr)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	feedPtr.unescapeText()
	return feedPtr, nil
}
