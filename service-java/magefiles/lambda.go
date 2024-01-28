//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Lambda mg.Namespace

// Start starts the lambda function locally
func (Lambda) Start() error {
	mg.Deps(Lambda.Build)
	mg.Deps(Cdk.Synth)

	return sh.RunV(
		"sam", "local", "start-api",
		"--template", "./deployment-aws/cdk.out/TaskResourceLocal.template.json",
		"--docker-network", "lambda-network",
		"--warm-containers", "EAGER",
	)
}

// Build builds all lambda functions
func (Lambda) Build() error {
	return sh.RunV("mvn", "package")
}
