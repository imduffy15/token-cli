// +build !windows

package config

import (
	"os"
)

func userHomeDir() string {
	return os.Getenv("HOME")
}

func makeDirectory() error {
	return os.MkdirAll(Dir(), 0755)
}
