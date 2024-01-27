package main

import (
	"root/go/lambda-v1-delete-task/incoming"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	handler, err := incoming.NewHandler()
	if err != nil {
		panic(err)
	}
	lambda.Start(handler.Invoke)
}
