package main

import (
	"os"
	"github.com/magefile/mage/sh"
)

// Clean cleans generated files
func Clean() error {
	paths := []string{
		"./deployment-aws/cdk.out",
	}

	for _, path := range paths {
		err := os.RemoveAll(path)
		if err != nil {
			return err
		}
	}

	return sh.RunV("mvn", "clean")
}
