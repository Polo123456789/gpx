package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"
)

var (
	DoneCheckbox = "\x1B[32m✓\x1B[0m"
)

func InstallPackage(ctx context.Context, p Package, binPath string) error {
	cmd := exec.CommandContext(ctx, "go", "install", p.String())
	cmd.Env = append(os.Environ(), "GOBIN="+binPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Print("Installing ", p.String(), "... ")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install %s: %w", p.String(), err)
	}

	fmt.Println(DoneCheckbox)

	if err := AddVersionToExecutable(p, binPath); err != nil {
		return err
	}

	return nil
}

func AddVersionToExecutable(p Package, binPath string) error {
	oldPath := path.Join(binPath, p.CommandName())
	newPath := path.Join(binPath, p.BinName())

	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to rename %s to %s: %w", oldPath, newPath, err)
	}
	return nil
}

func PackageIsInstalled(p Package, binPath string) bool {
	path := path.Join(binPath, p.BinName())
	_, err := os.Stat(path)
	return err == nil
}

// We can use the modification time to know when was the last time that we used
// it, and deleting it in a future 'clean' command

func TouchExecutable(p Package, binPath string) error {
	path := path.Join(binPath, p.BinName())
	now := time.Now().Local()
	if err := os.Chtimes(path, now, now); err != nil {
		return fmt.Errorf("failed to touch %s: %w", path, err)
	}
	return nil
}

func ExecPackage(
	ctx context.Context,
	p Package,
	binPath string,
	args ...string,
) error {
	cmd := exec.CommandContext(ctx, path.Join(binPath, p.BinName()), args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run %s: %w", p.String(), err)
	}

	return TouchExecutable(p, binPath)
}

type InstallAllCommand struct {
	forced   bool
	packages []Package
	binPath  string
}

var _ Command = &InstallAllCommand{}

func (c *InstallAllCommand) Name() string {
	return "i:install"
}

func (c *InstallAllCommand) Synopsis() string {
	return "install all packages"
}

func (c *InstallAllCommand) ParseFlags(args []string) {
	fset := flag.NewFlagSet(c.Name(), flag.ExitOnError)
	fset.BoolVar(&c.forced, "f", false, "force installation")
	_ = fset.Parse(args)
}

func (c *InstallAllCommand) Run(ctx context.Context, args []string) error {
	c.ParseFlags(args)

	for _, p := range c.packages {
		if !c.forced && PackageIsInstalled(p, c.binPath) {
			fmt.Printf("%s is already installed\n", p.String())
			continue
		}

		if err := InstallPackage(ctx, p, c.binPath); err != nil {
			return err
		}
	}
	return nil
}
