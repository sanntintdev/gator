package commands

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/sanntintdev/gator/internal/database"
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

func FetchFeed(ctx context.Context, url string) (*RSSFeed, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var feed RSSFeed
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return &feed, nil
}

func ScrapeFeeds(s *State) error {
	ctx := context.Background()
	feed, err := s.Db.RetrieveNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("couldn't get next feed: %w", err)
	}

	err = s.Db.MarkFeedFetched(ctx, feed.ID)
	if err != nil {
		return fmt.Errorf("couldn't mark feed as fetched: %w", err)
	}

	fmt.Printf("Fetching feed: %s from %s\n", feed.Name, feed.Url)
	rssFeed, err := FetchFeed(ctx, feed.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}

	fmt.Printf("\n=== %s ===\n", rssFeed.Channel.Title)
	fmt.Printf("Found %d posts:\n\n", len(rssFeed.Channel.Item))

	for i, item := range rssFeed.Channel.Item {
		fmt.Printf("%d. %s\n", i+1, item.Title)
		fmt.Printf("   Link: %s\n", item.Link)
		if item.PubDate != "" {
			fmt.Printf("   Published: %s\n", item.PubDate)
		}
		fmt.Println()
	}

	return nil
}

func handlerAgg(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Invalid number of arguments")
	}

	timeBetweenRequest, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Invalid time format: %w", err)
	}

	if timeBetweenRequest < 1*time.Second {
		return fmt.Errorf("time between requests must be at least 1s")
	}

	fmt.Printf("Collecting feeds every %s\n", timeBetweenRequest)
	fmt.Println("Press Ctrl+C to stop")

	tiker := time.NewTicker(timeBetweenRequest)
	defer tiker.Stop()

	for ; ; <-tiker.C {
		fmt.Println("Fetching feeds...")
		err := ScrapeFeeds(s)
		if err != nil {
			fmt.Printf("Error fetching feeds: %v\n", err)
		}
	}

}

func handlerCreateFeed(s *State, cmd Command, user database.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if len(cmd.Args) != 2 {
		return fmt.Errorf("Invalid number of arguments")
	}

	name := cmd.Args[0]
	feedUrl := cmd.Args[1]

	now := time.Now()
	createFeedParams := database.CreateFeedParams{
		Name:      name,
		Url:       feedUrl,
		UserID:    user.ID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	createdFeed, err := s.Db.CreateFeed(ctx, createFeedParams)
	if err != nil {
		return fmt.Errorf("Failed to create feed: %w", err)
	}

	// Created follow feed records
	followFeedParams := database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: createdFeed.ID,
	}

	_, err = s.Db.CreateFeedFollow(ctx, followFeedParams)
	if err != nil {
		return fmt.Errorf("Failed to create follow feed record: %w", err)
	}

	fmt.Printf("Feed created successfully.")
	return nil
}

func handlerRetrieveFeeds(s *State, cmd Command) error {
	ctx := context.Background()
	feeds, err := s.Db.RetrieveFeedsWithUser(ctx)

	if err != nil {
		return fmt.Errorf("Failed to retrieve feeds: %w", err)
	}

	fmt.Println("=== FEEDS ===")
	for _, feed := range feeds {
		fmt.Printf("  Name: %s\n", feed.Name)
		fmt.Printf("  URL: %s\n", feed.Url)
		fmt.Printf(" Created by: %s\n", feed.Name_2.String)
	}

	return nil
}

func handlerFollowFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Invalid number of arguments")
	}

	feedURL := cmd.Args[0]
	ctx := context.Background()
	feed, err := s.Db.RetrieveFeedWithURL(ctx, feedURL)

	if err != nil {
		return fmt.Errorf("Invalid feed URL: %w", err)
	}

	followFeedParams := database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	feedFollow, err := s.Db.CreateFeedFollow(ctx, followFeedParams)
	if err != nil {
		return fmt.Errorf("Failed to follow feed: %w", err)
	}

	fmt.Printf("Feed %s followed successfully.\n", feedFollow.FeedName)
	fmt.Printf("Followed by %s.\n", feedFollow.UserName)

	return nil
}

func handlerUnfollowFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Invalid number of arguments")
	}

	feedURL := cmd.Args[0]
	ctx := context.Background()
	feed, err := s.Db.RetrieveFeedWithURL(ctx, feedURL)

	if err != nil {
		return fmt.Errorf("Invalid feed URL: %w", err)
	}

	unfollowFeedParams := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	err = s.Db.DeleteFeedFollow(ctx, unfollowFeedParams)
	if err != nil {
		return fmt.Errorf("Failed to unfollow feed: %w", err)
	}

	fmt.Printf("Feed %s unfollowed successfully.\n", feed.Name)

	return nil
}

func RegisterFeedCommands(c *Commands) {
	publicHandlers := map[string]func(*State, Command) error{
		"agg":   handlerAgg,
		"feeds": handlerRetrieveFeeds,
	}

	authHandlers := map[string]func(*State, Command, database.User) error{
		"addfeed":  handlerCreateFeed,
		"follow":   handlerFollowFeed,
		"unfollow": handlerUnfollowFeed,
	}

	for name, handler := range publicHandlers {
		c.register(name, handler)
	}

	for name, handler := range authHandlers {
		c.register(name, MiddlewareLoggedIn(handler))
	}
}
