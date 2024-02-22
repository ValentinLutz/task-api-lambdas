//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Lambda mg.Namespace

// Build builds all lambda functions
func (Lambda) Build() error {
	return sh.RunV("mvn", "package")
}
