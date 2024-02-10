package main

import (
	"os"
)

// Clean cleans generated files
func Clean() error {
	paths := []string{
		"./lambda/bootstrap",
	}

	for _, path := range paths {
		err := os.RemoveAll(path)
		if err != nil {
			return err
		}
	}

	return nil
}
