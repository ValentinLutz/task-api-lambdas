//go:build mage
// +build mage

package main

import (
	"os"

	library "root/library-golang"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Docker mg.Namespace

// Up starts the docker-compose stack
func (Docker) Up() error {
	library.GetOrSetDefaultDatabaseEnvVars()

	os.Chdir("../deployment-docker")
	defer os.Chdir("..")

	return sh.RunV(
		"docker",
		"compose",
		"--file", "./docker-compose.yaml",
		"up",
		"--force-recreate",
		//"--detach",
		"--wait",
	)
}

// Down stops the docker-compose stack
func (Docker) Down() error {
	library.GetOrSetDefaultDatabaseEnvVars()

	os.Chdir("../deployment-docker")
	defer os.Chdir("..")

	return sh.RunV(
		"docker",
		"compose",
		"--file", "./docker-compose.yaml",
		"down",
	)
}

// Logs shows the logs of the docker-compose stack
func (Docker) Logs() error {
	library.GetOrSetDefaultDatabaseEnvVars()

	os.Chdir("../deployment-docker")
	defer os.Chdir("..")

	return sh.RunV(
		"docker",
		"compose",
		"--file", "./docker-compose.yaml",
		"logs",
	)
}
