package main

import (
	"boot-dev-gator/internal/database"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func handlerLogin(s *State, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("login handler expects a single argument, the username")
	}

	user, err := s.db.GetUserByName(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Println(user)

	err = s.config.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Success! Current user set to %v", cmd.args[0])
	return nil
}

func handlerRegister(s *State, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("register handler expects a single argument, the username")
	}

	now := time.Now()
	args := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      cmd.args[0],
	}
	user, err := s.db.CreateUser(context.Background(), args)
	if err != nil {
		return err
	}

	s.config.SetUser(user.Name)
	fmt.Printf("User %v was created", user.Name)
	fmt.Println(user)

	return nil
}

func handlerReset(s *State, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("reset handler requires no arguments")
	}

	if err := s.db.DeleteAllUsers(context.Background()); err != nil {
		return err
	}

	if err := s.db.DeleteAllFeeds(context.Background()); err != nil {
		return err
	}

	if err := s.db.DeleteAllFeedFollows(context.Background()); err != nil {
		return err
	}

	return nil
}

func handlerUsers(s *State, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("user handler requires no arguments")
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Name == s.config.CurrentUserName {
			fmt.Println(user.Name + " (current)")
		} else {
			fmt.Println(user.Name)
		}
	}

	return nil
}

func handlerAgg(s *State, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("agg handler requires one argument, time_between_reqs")
	}

	duration, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %vm%vs\n", duration.Minutes(), duration.Seconds())

	ticker := time.NewTicker(duration)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func handlerAddFeed(s *State, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return errors.New("addfeed handler requires two arguments, feed name and feed url")
	}

	now := time.Now()
	args := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), args)
	if err != nil {
		return err
	}

	argss := database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	s.db.CreateFeedFollows(context.Background(), argss)
	return nil
}

func handlerAllFeeds(s *State, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("feeds handler requires no arguments")
	}

	feeds, err := s.db.GetAllFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("Name: %v URL: %v User: %v\n", feed.Name, feed.Url, user.Name)
	}

	return nil
}

func handlerFollow(s *State, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("follow handler expects a single argument, the url")
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	now := time.Now()
	args := database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	feedFollow, err := s.db.CreateFeedFollows(context.Background(), args)
	if err != nil {
		return err
	}

	fmt.Println(feedFollow.FeedName, user)
	return nil
}

func handlerUnfollow(s *State, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("follow handler expects a single argument, the url")
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	args := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	err = s.db.DeleteFeedFollow(context.Background(), args)
	if err != nil {
		return err
	}

	return nil
}

func handlerFollowing(s *State, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return errors.New("following handler expects no arguments")
	}

	feedFollows, err := s.db.GetAllFeedFollowsByUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, feedFollow := range feedFollows {
		fmt.Println(feedFollow.FeedName)
	}

	return nil
}

func handlerBrowse(s *State, cmd command, user database.User) error {
	var limit int
	if len(cmd.args) == 0 {
		limit = 2
	}
	if len(cmd.args) == 1 {
		val, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return err
		}
		limit = val
	}

	fmt.Printf("Getting %v posts", limit)
	args := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}
	posts, err := s.db.GetPostsForUser(context.Background(), args)
	if err != nil {
		return err
	}
	fmt.Printf("%v posts found", len(posts))
	for i, post := range posts {
		fmt.Printf("%v: %v", i, post.Title)
	}

	return nil

}
