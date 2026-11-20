package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/Cyrus-Gahatraj/gator/internal/database"
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


func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("get current user: %v\n", err)
		}

		return handler(s, cmd, user)
	}
}


func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}
	req.Header.Set("User-Agent", "gator")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, err
	}

	var rssFeed RSSFeed
	if err := xml.Unmarshal(data, &rssFeed); err != nil {
		return &RSSFeed{}, err
	}

	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)

	for i := range rssFeed.Channel.Item {
		rssFeed.Channel.Item[i].Description = html.UnescapeString(rssFeed.Channel.Item[i].Description)
		rssFeed.Channel.Item[i].Title = html.UnescapeString(rssFeed.Channel.Item[i].Title)
	}

	return &rssFeed, nil
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("could not get next feed: %w", err)
	}

	err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("could not mark feed as fetched: %w", err)
	}

	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("could not fetch feed %s: %w", feed.Name, err)
	}

	for _, item := range feedData.Channel.Item {
		fmt.Printf("* %s\n", item.Title)
	}

	return nil
}

func aggregator(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <time_between_reqs>", cmd.name)
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	fmt.Printf("Collecting feeds every %s\n", cmd.args[0])

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			fmt.Printf("error: %v\n", err)
		}
	}
}

func addFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("addfeed need name and url as arguments")
	}

	currentTime := time.Now()
	feedParam := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), feedParam)
	if err != nil {
		return fmt.Errorf("%v\n", err)
	}

	feedFollowParam := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), feedFollowParam)
	if err != nil {
		return fmt.Errorf("failed to create feed follow: %w", err)
	}

	fmt.Printf("%+v\n", feed)
	printFeedFollow(feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func feeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("feed need no argument")
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("%v\n", err)
	}

	for _, feed := range feeds {
		user, err := s.db.GetUserFromID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("Can't fetch username of the post %v: %v", feed.Url, err)
		}
		fmt.Printf("Name: %v\nURL: %v\nUsername: %v\n\n", feed.Name, feed.Url, user.Name)
	}

	return nil
}

func follow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("follow need a single url as argument")
	}

	feed, err := s.db.GetFeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("could not get feed: %w", err)
	}

	currentTime := time.Now()
	feedFollowParam := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), feedFollowParam)
	if err != nil {
		return fmt.Errorf("could not create feed follow: %w", err)
	}

	printFeedFollow(feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func following(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("following need no argument")
	}

	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("could not get feed follows: %w", err)
	}

	for _, ff := range feedFollows {
		fmt.Printf("* %s\n", ff.FeedName)
	}

	return nil
}

func unfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("unfollow need a single url as agrument")
	}

	feed, err := s.db.GetFeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("could not get feed: %w", err)
	}

	deleteParam := database.DeleteFeedFollowParams {
		UserID: user.ID,
		FeedID: feed.ID,
	}

	if err := s.db.DeleteFeedFollow(context.Background(), deleteParam); err != nil {
		return fmt.Errorf("could not create feed follow: %w", err)
	}

	return nil
}

func printFeedFollow(username, feedname string) {
	fmt.Printf("* User:          %s\n", username)
	fmt.Printf("* Feed:          %s\n", feedname)
}
