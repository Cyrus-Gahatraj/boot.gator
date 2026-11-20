package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Cyrus-Gahatraj/boot.gator/internal/database"
)

func browse(s *state, cmd command, user database.User) error {
	limit := 2
	switch len(cmd.args) {
	case 0:
	case 1:
		parsed, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("limit must be a number: %w", err)
		}
		if parsed < 1 {
			return fmt.Errorf("limit must be at least 1")
		}
		limit = parsed
	default:
		return fmt.Errorf("usage: %s [limit]", cmd.name)
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("could not get posts: %w", err)
	}

	for _, post := range posts {
		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Printf("Published: %s\n", post.PublishedAt.Format("2006-01-02 15:04:05 MST"))
		fmt.Printf("Description: %s\n", post.Description.String)
		fmt.Println()
	}

	return nil
}
