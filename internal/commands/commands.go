package commands

import (
	"fmt"

	"github.com/sanntintdev/gator/internal/config"
	"github.com/sanntintdev/gator/internal/database"
)

type State struct {
	Db  *database.Queries
	Cfg *config.Config
}

type Command struct {
	Name string
	Args []string
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
	RegisterUserCommands(c)
}
