// main.go
package main

import (
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/robbyriverside/fibber"
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

	_, err := parser.Parse()
	if err != nil {
		os.Exit(0)
	}
}
