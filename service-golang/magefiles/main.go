package main

import (
	"os"
)

// Clean cleans generated files
func Clean() error {
	paths := []string{
		"./deployment-aws/cdk.out",
		"./lambda-v1-delete-task/incoming/model.gen.go",
		"./lambda-v1-get-task/incoming/model.gen.go",
		"./lambda-v1-get-tasks/incoming/model.gen.go",
		"./lambda-v1-post-tasks/incoming/model.gen.go",
		"./lambda-v1-put-task/incoming/model.gen.go",
		"./lambda-v1-delete-task/bootstrap",
		"./lambda-v1-get-task/bootstrap",
		"./lambda-v1-get-tasks/bootstrap",
		"./lambda-v1-post-tasks/bootstrap",
		"./lambda-v1-put-task/bootstrap",
	}

	for _, path := range paths {
		err := os.RemoveAll(path)
		if err != nil {
			return err
		}
	}

	return nil
}
