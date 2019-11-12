package main

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/url"
	"slack-lottery-bot/adaptor"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nlopes/slack"
)

type handler struct {
	verificationToken string
	api               adaptor.API
}

type Handler interface {
	Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

func NewHandler(verificationToken string, botToken string) Handler {
	return &handler{verificationToken, adaptor.NewAPI(botToken)}
}

func (h *handler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	reqBody := request.Body

	message := new(slack.InteractionCallback)
	jsonStr, err := url.QueryUnescape(reqBody[8:])
	if err != nil {
		log.Printf("failed to unespace request body: %s", err)
		return events.APIGatewayProxyResponse{}, err
	}

	if err := json.Unmarshal([]byte(jsonStr), message); err != nil {
		log.Printf("failed to decode json message from slack: %s", jsonStr)
		return events.APIGatewayProxyResponse{}, err
	}

	if message.Token != h.verificationToken {
		log.Printf("invalid token: %s", message.Token)
		return events.APIGatewayProxyResponse{}, errors.New("invalid token")
	}

	action := message.ActionCallback.AttachmentActions[0]
	switch action.Name {
	case "select":
		log.Print("select action")
		err := h.lottery(action.Value, message.Channel.ID)
		if err != nil {
			log.Print(err)
			return events.APIGatewayProxyResponse{}, err
		}

		originalMessage := message.OriginalMessage
		jsonBody, err := responseMessage(&originalMessage, ":ok:", "")
		if err != nil {
			log.Print(err)
			return events.APIGatewayProxyResponse{}, err
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(jsonBody),
		}, nil

	case "cancel":
		originalMessage := message.OriginalMessage
		jsonBody, err := responseMessage(&originalMessage, "キャンセルしました", "")
		if err != nil {
			log.Print(err)
			return events.APIGatewayProxyResponse{}, err
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(jsonBody),
		}, nil

	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Bad Request",
		}, nil
	}
}

func (h *handler) lottery(actionValue string, channelID string) error {
	var userIDs []string
	var err error

	if actionValue == "channel" {
		userIDs, err = h.api.GetUsersFromChannel(channelID)
	} else {
		userIDs, err = h.api.GetUsersFromUserGroup(actionValue)
	}

	if err != nil {
		log.Print(err)
		return err
	}

	userID := lotteryOneUserFromUsers(userIDs)
	return h.api.PostMessage(channelID, "<@"+userID+"> が当選しました")
}

func responseMessage(original *slack.Message, titie, value string) ([]byte, error) {
	original.ReplaceOriginal = true
	original.Attachments[0].Actions = []slack.AttachmentAction{}
	original.Attachments[0].Fields = []slack.AttachmentField{
		{
			Title: titie,
			Value: value,
			Short: false,
		},
	}
	jsonBody, err := json.Marshal(original)
	return jsonBody, err
}

func lotteryOneUserFromUsers(userIDs []string) string {
	rand.Seed(time.Now().UnixNano())
	userID := userIDs[rand.Intn(len(userIDs)-1)]
	return userID
}
