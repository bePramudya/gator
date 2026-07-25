package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bePramudya/gator/internal/database"
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
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var feeds = &RSSFeed{}
	if err = xml.Unmarshal(data, feeds); err != nil {
		return nil, err
	}

	feeds.Channel.Title = html.UnescapeString(feeds.Channel.Title)
	feeds.Channel.Description = html.UnescapeString(feeds.Channel.Description)

	for i, item := range feeds.Channel.Item {
		feeds.Channel.Item[i].Title = html.UnescapeString(item.Title)
		feeds.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}

	return feeds, nil
}

func handlerAggregate(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("invalid time: agg <time>")
	}

	fmt.Printf("Collecting feeds every %v \n", cmd.args[0])

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		fmt.Println("fetching batch...")
		scrapeFeeds(s)
		fmt.Println("Finished fetching")
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("input error: addfeed <name> <url>")
	}

	createFeedArgs := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), createFeedArgs)
	if err != nil {
		fmt.Println("error creating feed")
		return err
	}

	createFeedFollowArgs := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    feed.UserID,
		FeedID:    feed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), createFeedFollowArgs)
	if err != nil {
		return err
	}

	fmt.Println(feed)
	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("name: %v\n", feed.Name)
		fmt.Printf(" * url: %v\n", feed.Url)
		fmt.Printf(" * creator: %v\n", feed.UserName.String)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("provide one url: follow <url>")
	}

	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	args := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), args)
	if err != nil {
		return err
	}

	fmt.Printf("Feed: %v", feedFollow.FeedName)
	fmt.Printf("Name: %v", feedFollow.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	following, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, follow := range following {
		fmt.Println(follow.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	args := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	if err = s.db.DeleteFeedFollow(context.Background(), args); err != nil {
		return err
	}

	return nil
}

func handlerBrowse(s *state, cmd command) error {
	var limit int32 = 2
	if len(cmd.args) == 1 {
		lmt, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return err
		}

		limit = int32(lmt)
	}

	posts, err := s.db.GetPostForUser(context.Background(), limit)
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Printf("# %v\n", post.Title)
		fmt.Printf(" * URL: %v\n", post.Url)
		fmt.Printf(" * Description: %v\n", post.Description)
		fmt.Printf(" * Published Date: %v\n\n\n", post.PublishedAt)
	}

	return nil
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	markFeedArgs := database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: time.Now(),
		ID:        feed.ID,
	}

	feed, err = s.db.MarkFeedFetched(context.Background(), markFeedArgs)
	if err != nil {
		return err
	}

	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	for _, rf := range rssFeed.Channel.Item {
		publishedAt := sql.NullTime{}
		parsed, err := time.Parse(time.RFC1123Z, rf.PubDate)
		if err == nil {
			publishedAt = sql.NullTime{
				Time:  parsed,
				Valid: true,
			}
		}

		createPostArgs := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       rf.Title,
			Url:         rf.Link,
			Description: rf.Description,
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		}

		_, err = s.db.CreatePost(context.Background(), createPostArgs)
		if err != nil {
			if strings.Contains(err.Error(), "unique") {
				continue
			} else {
				return err
			}
		}

		fmt.Println(rf.Title)
	}

	return nil
}
