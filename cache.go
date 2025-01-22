package main

import (
	"fmt"
	"os"
	"path/filepath"
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
