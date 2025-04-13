// main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jessevdk/go-flags"
	"github.com/robbyriverside/fibber"
	config "github.com/robbyriverside/fibber/config"
)

var (
	Version   string // set via -ldflags
	Commit    string
	BuildTime string
)

type InsertCommand struct {
	ID      string `short:"i" long:"id" description:"ID of the gem" required:"true"`
	Content string `short:"c" long:"content" description:"Raw content of the gem" required:"true"`
}

type StatusCommand struct{}

type InspectCommand struct{}

type VersionCommand struct{}
type ConfigCommand struct{}

type ConfigSetCommand struct {
	Args struct {
		Key   string `positional-arg-name:"key" required:"true"`
		Value string `positional-arg-name:"value" required:"true"`
	} `positional-args:"yes"`
}

type ConfigGetCommand struct {
	Args struct {
		Key string `positional-arg-name:"key" required:"true"`
	} `positional-args:"yes"`
}

type ConfigDescribeCommand struct{}

func (cmd *ConfigDescribeCommand) Execute(args []string) error {
	fmt.Println("Fibber config file:", config.Path())

	lines, err := config.Describe()
	if err != nil {
		return err
	}
	for _, line := range lines {
		fmt.Println(line)
	}
	return nil
}

func (cmd *VersionCommand) Execute(args []string) error {
	fmt.Printf("%s version %s\n", fibber.Options.AppName, fibber.Options.Version)
	if fibber.Options.Commit != "" {
		fmt.Printf("Commit:    %s\n", fibber.Options.Commit)
	}
	if fibber.Options.BuildTime != "" {
		fmt.Printf("Built at:  %s\n", fibber.Options.BuildTime)
	}
	return nil
}

func (cmd *InsertCommand) Execute(args []string) error {
	fibber.InsertGem(cmd.ID, cmd.Content)
	return nil
}

func (cmd *StatusCommand) Execute(args []string) error {
	fibber.StatusReport()
	return nil
}

func (cmd *InspectCommand) Execute(args []string) error {
	fibber.Inspect()
	return nil
}

func (cmd *ConfigCommand) Execute(args []string) error {
	fmt.Println("Fibber config file:", config.Path())

	// Ensure the config file exists
	if _, err := os.Stat(config.Path()); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(config.Path()), 0755); err != nil {
			return fmt.Errorf("unable to create config directory: %w", err)
		}
		if _, err := os.Create(config.Path()); err != nil {
			return fmt.Errorf("unable to touch config file: %w", err)
		}
	}

	// Show config keys and their current/default values
	fmt.Println()
	fmt.Println("Available configuration keys:")

	lines, err := config.Describe()
	if err != nil {
		return fmt.Errorf("failed to describe config: %w", err)
	}
	for _, line := range lines {
		fmt.Println(line)
	}

	fmt.Println()
	fmt.Println("Use:")
	fmt.Println("  fibber config get <key>")
	fmt.Println("  fibber config set <key> <value>")
	fmt.Println("  fibber config init       (overwrite defaults)")

	return nil
}

func (cmd *ConfigGetCommand) Execute(args []string) error {
	value, err := config.Get(cmd.Args.Key)
	if err != nil {
		return err
	}
	fmt.Println(value)
	return nil
}

func (cmd *ConfigSetCommand) Execute(args []string) error {
	if err := config.Set(cmd.Args.Key, cmd.Args.Value); err != nil {
		return fmt.Errorf("error setting config key '%s': %w", cmd.Args.Key, err)
	}
	value, _ := config.Get(cmd.Args.Key)
	fmt.Printf("%s = %s\n", cmd.Args.Key, value)
	return nil
}

func main() {
	fibber.Options.AppName = "fibber"
	fibber.Options.Version = Version
	fibber.Options.Commit = Commit
	fibber.Options.BuildTime = BuildTime
	fibber.InitLogger(os.Getenv("ENV"))

	parser := flags.NewParser(fibber.Options, flags.Default)
	parser.AddCommand("insert", "Insert a gem", "Adds a new gem and updates the composition", &InsertCommand{})
	parser.AddCommand("status", "Composition status", "Reports recent activity and suggested next steps", &StatusCommand{})
	parser.AddCommand("inspect", "Inspect model", "Prints internal model for debugging", &InspectCommand{})
	parser.AddCommand("version", "Show version", "Displays app version, commit, and build time", &VersionCommand{})
	configCmd := &ConfigCommand{}
	configParser, _ := parser.AddCommand("config", "Manage config", "Inspect or modify configuration", configCmd)

	configParser.AddCommand("describe", "Show config", "Show current configuration", &ConfigDescribeCommand{})
	configParser.AddCommand("set", "Set config key", "Set a specific config key", &ConfigSetCommand{})
	configParser.AddCommand("get", "Get config key", "Retrieve a config value", &ConfigGetCommand{})

	_, err := parser.Parse()
	if err != nil {
		os.Exit(0)
	}
}
