package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
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
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error: could not perform request \n%v", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error: could not read response body")
	}
	feed := RSSFeed{}
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, fmt.Errorf("error: could not unmarshal http response unto RSSFeed struct \n%v", err)
	}
	feed.unescapeRSSFeed()
	return &feed, err
}

func (f *RSSFeed) unescapeRSSFeed() {
	f.Channel.Description = html.UnescapeString(f.Channel.Description)
	f.Channel.Title = html.UnescapeString(f.Channel.Title)
	for i := range f.Channel.Item {
		f.Channel.Item[i].Description = html.UnescapeString(f.Channel.Item[i].Description)
		f.Channel.Item[i].Title = html.UnescapeString(f.Channel.Item[i].Title)
	}
}
