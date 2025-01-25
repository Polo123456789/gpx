package main

import (
	"context"
	"errors"
)

type Command interface {
	Run(ctx context.Context, args []string) error
	Name() string
	Synopsis() string
}

func RunCommand(ctx context.Context, commands []Command, args []string) error {
	if len(args) < 2 {
		return errors.New(GenerateHelp(commands))
	}

	command := args[1]
	arguments := []string{}
	if len(args) > 2 {
		arguments = args[2:]
	}

	for _, c := range commands {
		if c.Name() == command {
			return c.Run(ctx, arguments)
		}
	}

	return errors.New("Command not found\n\n" + GenerateHelp(commands))
}

func GenerateHelp(commands []Command) string {

	help := `Golang Package eXecute

Run commands from your tools.go always using the same version that you have in
your go.mod file.

Available commands:

`
	for _, c := range commands {
		help += "\t" + c.Name() + " - " + c.Synopsis() + "\n"
	}

	return help
}
