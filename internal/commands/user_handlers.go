package commands

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sanntintdev/gator/internal/database"
)

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("Login command requires an username")
	}

	username := cmd.Args[0]

	ctx := context.Background()
	_, err := s.Db.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	err = s.Cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}
	fmt.Printf("User successfully set to: %s\n", username)
	return nil
}

func handlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("Register command requires an username")
	}

	username := cmd.Args[0]
	ctx := context.Background()

	_, err := s.Db.GetUser(ctx, username)
	if err != nil {
		fmt.Printf("Error: User '%s' already exists\n", username)
	}

	now := time.Now()
	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		Name:      username,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err = s.Db.CreateUser(ctx, userParams)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	err = s.Cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("failed to set current user: %w", err)
	}

	fmt.Printf("User successfully created: %s\n", username)

	return nil
}

func handlerReset(s *State, c Command) error {
	err := s.Db.ResetAllUser(context.Background())
	if err != nil {
		return fmt.Errorf("failed to reset users: %w", err)
	}
	fmt.Println("All users successfully reset")
	return nil
}

func handlerGetUsers(s *State, c Command) error {
	ctx := context.Background()
	users, err := s.Db.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	for _, user := range users {
		if user.Name == s.Cfg.CurrentUserName {
			fmt.Println("*", user.Name, "(current)")
			continue
		}
		fmt.Println("*", user.Name)
	}
	return nil
}

func handlerFollowing(s *State, cmd Command, user database.User) error {

	ctx := context.Background()

	following, err := s.Db.RetrieveFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to get following: %w", err)
	}

	for _, following := range following {
		fmt.Println("*", following.FeedName)
	}
	return nil
}

func RegisterUserCommands(c *Commands) {
	userHandlers := map[string]func(*State, Command) error{
		"login":    handlerLogin,
		"register": handlerRegister,
		"reset":    handlerReset,
		"users":    handlerGetUsers,
	}

	for name, handler := range userHandlers {
		c.register(name, handler)
	}

	c.register("following", MiddlewareLoggedIn(handlerFollowing))
}
