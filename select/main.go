package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	verificationToken := os.Getenv("VERIFICATION_TOKEN")
	botToken := os.Getenv("BOT_TOKEN")
	oauthToken := os.Getenv("OAUTH_TOKEN")

	handler := NewHandler(verificationToken, botToken, oauthToken)
	lambda.Start(handler.Handle)
}