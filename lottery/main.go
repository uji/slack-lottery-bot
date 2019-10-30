package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	handler := NewHandler(token)
	lambda.Start(handler.Handle)
}
