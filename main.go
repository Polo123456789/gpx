package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	toolsSrc, err := os.ReadFile("tools.go")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading tools.go:", err)
		fmt.Println("Remember to run this command in the root of your project")
		return
	}

	modSrc, err := os.ReadFile("go.mod")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading go.mod:", err)
		fmt.Println("Remember to run this command in the root of your project")
		return
	}

	binPath, err := GetBinariesPath()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting cache path:", err)
		return
	}

	packages, err := ListPackages(string(toolsSrc), string(modSrc))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error listing packages:", err)
		return
	}

	commands := []Command{}

	commands = append(
		commands,
		&InstallAllCommand{
			binPath:  binPath,
			packages: packages,
		},
		&CleanCacheCommand{
			binPath: binPath,
		},
	)

	for _, p := range packages {
		commands = append(commands, &ExecCommand{
			pkg:     p,
			binPath: binPath,
		})
	}

	if err := RunCommand(ctx, commands, os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func GetBinariesPath() (string, error) {
	binPath, err := GetCachePath()
	if err != nil {
		return "", err
	}

	err = CheckCachePath(binPath)
	if err != nil {
		return "", err
	}
	return binPath, nil
}
