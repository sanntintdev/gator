package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sanntintdev/gator/internal/commands"
	"github.com/sanntintdev/gator/internal/config"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	appState := &commands.State{
		Cfg: &cfg,
	}

	cmds := commands.NewCommands()
	cmds.RegisterDefaultCommands()

	err = parsedArgsAndExecute(cmds, appState)

	if err != nil {
		log.Fatalf("Failed to execute command: %v", err)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func parsedArgsAndExecute(cmds *commands.Commands, state *commands.State) error {
	args := os.Args
	if len(args) < 2 {
		return fmt.Errorf("No command provided")
	}

	cmd := commands.Command{
		Name: args[1],
		Args: []string{},
	}

	if len(args) > 2 {
		cmd.Args = args[2:]
	}

	return cmds.Run(state, cmd)
}
