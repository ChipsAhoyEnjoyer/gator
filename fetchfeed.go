package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

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
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
	}
	return &feed, err
}
