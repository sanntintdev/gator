package commands

import (
	"context"
	"fmt"

	"github.com/sanntintdev/gator/internal/database"
)

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		username := s.Cfg.CurrentUserName
		ctx := context.Background()

		currentUser, err := s.Db.GetUser(ctx, username)
		if err != nil {
			return fmt.Errorf("failed to get user %s: %w", username, err)
		}

		return handler(s, cmd, currentUser)
	}
}
