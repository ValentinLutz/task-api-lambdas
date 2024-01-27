//go:build mage
// +build mage

package main

import (
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Test mg.Namespace

// Functional runs the integration tests
func (Test) Functional() error {
	os.Chdir("./test-functional")
	defer os.Chdir("..")

	return sh.RunV("go", "test", "-count=1", "-p=1", "./...")
}
