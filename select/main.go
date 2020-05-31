package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	verificationToken := os.Getenv("VERIFICATIONTOKEN")
	botToken := os.Getenv("BOTTOKEN")
	oauthToken := os.Getenv("OAUTHTOKEN")

	handler := NewHandler(verificationToken, botToken, oauthToken)
	lambda.Start(handler.Handle)
}
