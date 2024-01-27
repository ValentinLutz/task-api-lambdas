package testfunctional

import "os"

const BaseUrl = "http://localhost:8080/v1/tasks"

func ReadFile(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	return file
}
