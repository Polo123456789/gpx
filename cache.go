package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const CacheSubdirectory = "gpx"

func GetCachePath() (string, error) {
	p, ok := os.LookupEnv("GPX_BIN")
	if ok {
		return p, nil
	}

	var err error
	p, err = os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("could not determine user cache dir, mabye set $GPX_BIN? error: %v", err)
	}

	return filepath.Join(p, CacheSubdirectory), nil
}

// CheckCachePath ensures that the directory exists, and that we can write to it
func CheckCachePath(p string) error {
	exists, err := CacheExists(p)
	if err != nil {
		return err
	}

	if !exists {
		err = os.Mkdir(p, 0755)
		if err != nil {
			return fmt.Errorf("cold not create gpx cache dir in %s: %v", p, err)
		}
	}

	return nil
}

func CacheExists(p string) (bool, error) {
	info, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	if !info.IsDir() {
		return false, fmt.Errorf("%s exists, but is not a directory", p)
	}

	return true, nil
}

func GetFileAge(p string) (time.Duration, error) {
	info, err := os.Stat(p)
	if err != nil {
		return 0, err
	}

	return time.Since(info.ModTime()), nil
}

type CleanCacheCommand struct {
	binPath   string
	olderThan time.Duration
}

func (c *CleanCacheCommand) Name() string {
	return "i:clean"
}

func (c *CleanCacheCommand) Synopsis() string {
	return "clean old binaries from cache"
}

func (c *CleanCacheCommand) ParseFlags(args []string) {
	fset := flag.NewFlagSet(c.Name(), flag.ExitOnError)
	fset.DurationVar(
		&c.olderThan,
		"older-than",
		30*24*time.Hour,
		"delete files older than this duration",
	)
	_ = fset.Parse(args)
}

func (c *CleanCacheCommand) Run(ctx context.Context, args []string) error {
	c.ParseFlags(args)

	files, err := os.ReadDir(c.binPath)
	if err != nil {
		return fmt.Errorf("could not read cache dir %s: %v", c.binPath, err)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		age, err := GetFileAge(filepath.Join(c.binPath, f.Name()))
		if err != nil {
			return fmt.Errorf("could not get age of %s: %v", f.Name(), err)
		}

		if age > c.olderThan {
			fmt.Print("Removing ", f.Name(), " ... ")
			err = os.Remove(filepath.Join(c.binPath, f.Name()))
			if err != nil {
				return fmt.Errorf("could not remove %s: %v", f.Name(), err)
			}
			fmt.Println(DoneCheckbox)
		}
	}

	return nil
}
