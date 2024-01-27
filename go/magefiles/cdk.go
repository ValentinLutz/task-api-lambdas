//go:build mage
// +build mage

package main

import (
	"os"
	"root/library"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Cdk mg.Namespace

// Synth synthesizes the CDK stack
func (Cdk) Synth() error {
	library.GetOrSetDefaultStageEnvVars()

	os.Chdir("./deployment-aws")
	defer os.Chdir("..")

	return sh.RunV(
		"cdk",
		"synth",
	)
}

// Diff diffs the CDK stack
func (Cdk) Diff() error {
	library.GetOrSetDefaultStageEnvVars()

	os.Chdir("./deployment-aws")
	defer os.Chdir("..")

	return sh.RunV(
		"cdk",
		"diff",
	)
}

// Deploy deploys the CDK stack
func (Cdk) Deploy() error {
	library.GetOrSetDefaultStageEnvVars()

	os.Chdir("./deployment-aws")
	defer os.Chdir("..")

	return sh.RunV(
		"cdk",
		"deploy",
	)
}

// Destroy destroys the CDK stack
func (Cdk) Destroy() error {
	library.GetOrSetDefaultStageEnvVars()

	os.Chdir("./deployment-aws")
	defer os.Chdir("..")

	return sh.RunV(
		"cdk",
		"destroy",
	)
}
