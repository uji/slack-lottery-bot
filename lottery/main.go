package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nlopes/slack"
)

func main() {
	token := os.Getenv("BOT_TOKEN")
	client := slack.New(token)

	handler := NewHandler(client, token)
	lambda.Start(handler.Handle)
}
