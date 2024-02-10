//go:build mage
// +build mage

package main

import (
	"fmt"
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

// Load runs the load tests
func (Test) Load() error {
	os.Chdir("./test-load")
	defer os.Chdir("..")

	name := fmt.Sprintf("k6-test-go")

	return sh.RunV(
		"docker",
		"run",
		"--rm",
		"--name", name,
		"--volume", "./script.js:/k6/script.js:ro",
		"--volume", "./results:/results:rw",
		"--network", "host",
		"ghcr.io/grafana/xk6-dashboard:0.7.2",
		"run",
		"--out", fmt.Sprintf("web-dashboard=export=/results/%s.html", name),
		"--tag", fmt.Sprintf("testid=%s", name),
		"/k6/script.js",
	)
}
