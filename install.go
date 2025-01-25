package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"
)

func InstallPackage(ctx context.Context, p Package, binPath string) error {
	cmd := exec.CommandContext(ctx, "go", "install", p.String())
	cmd.Env = append(os.Environ(), "GOBIN="+binPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Installing package, command:", cmd.String())

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install %s: %w", p.String(), err)
	}

	AddVersionToExecutable(p, binPath)

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
	if err != nil {
		return false
	}
	return true
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

func ExecPackage(ctx context.Context, p Package, binPath string, args ...string) error {
	cmd := exec.CommandContext(ctx, path.Join(binPath, p.BinName()), args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run %s: %w", p.String(), err)
	}

	return TouchExecutable(p, binPath)
}
