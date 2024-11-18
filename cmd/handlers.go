package main

import (
	"boot-dev-gator/internal/database"
	"context"
	"errors"
	"fmt"
	"html"
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
	if len(cmd.args) != 0 {
		return errors.New("user handler requires no arguments")
	}

	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
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
		item.Link = item.Link

		cleanFeed.Channel.Item = append(cleanFeed.Channel.Item, item)
	}

	return nil
}

func handlerAddFeed(s *State, cmd command) error {
	if len(cmd.args) != 2 {
		return errors.New("addfeed handler requires two arguments, feed name and feed url")
	}
	user, err := s.db.GetUserByName(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return err
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

	_, err = s.db.CreateFeed(context.Background(), args)
	if err != nil {
		return err
	}

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
