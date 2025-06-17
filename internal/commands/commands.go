package commands

import (
	"errors"
	"fmt"

	"github.com/sanntintdev/gator/internal/config"
)

type State struct {
	Cfg *config.Config
}

type Command struct {
	Name string
	Args []string
}

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("Login command requires an username")
	}

	username := cmd.Args[0]
	err := s.Cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}
	fmt.Printf("User successfully set to: %s\n", username)
	return nil
}

type Commands struct {
	handler map[string]func(*State, Command) error
}

func NewCommands() *Commands {
	return &Commands{
		handler: make(map[string]func(*State, Command) error),
	}
}

func (c *Commands) Run(s *State, cmd Command) error {
	handler, ok := c.handler[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}
	return handler(s, cmd)
}

func (c *Commands) register(name string, handler func(*State, Command) error) {
	c.handler[name] = handler
}

func (c *Commands) RegisterDefaultCommands() {
	c.register("login", handlerLogin)
}
