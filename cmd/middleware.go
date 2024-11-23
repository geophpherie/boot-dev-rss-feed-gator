package main

import (
	"boot-dev-gator/internal/database"
	"context"
)

func middlewareLoggedIn(handler func(s *State, cmd command, user database.User) error) func(*State, command) error {
	return func(s *State, cmd command) error {
		user, err := s.db.GetUserByName(context.Background(), s.config.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}
