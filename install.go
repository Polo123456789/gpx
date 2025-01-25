package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
)

func InstallPackage(ctx context.Context, p Package, binPath string) error {
	cmd := exec.CommandContext(ctx, "go", "install", p.String())
	cmd.Env = append(os.Environ(), "GOBIN="+binPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Installing Package:", cmd.String())

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install %s: %w", p.String(), err)
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
