package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)

type handler struct {
	slackClient       *slack.Client
	verificationToken string
}

type Handler interface {
	Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

func NewHandler(client *slack.Client, token string) Handler {
	return &handler{client, token}
}

var (
	verificationToken = os.Getenv("VERIFICATION_TOKEN")
	api               = slack.New(os.Getenv("BOT_TOKEN"))
)

func (h *handler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	reqBody := request.Body
	eventsAPIEvent, err := slackevents.ParseEvent(
		json.RawMessage(reqBody),
		slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: verificationToken}),
	)
	if err != nil {
		log.Print(err)
		return events.APIGatewayProxyResponse{}, err
	}

	log.Print(eventsAPIEvent.Type)
	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(reqBody), &r)
		if err != nil {
			log.Print(err)
			return events.APIGatewayProxyResponse{}, err
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       r.Challenge,
		}, nil
	}

	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		log.Print(innerEvent.Type)
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			memberID, err := getOneUserFromChannel(ev.Channel)
			if err != nil {
				log.Print(err)
				return events.APIGatewayProxyResponse{}, err
			}
			api.PostMessage(ev.Channel, slack.MsgOptionText("<@"+memberID+"> が当選しました", false))
			return events.APIGatewayProxyResponse{
				StatusCode: 200,
			}, nil
		default:
			log.Print(ev)
		}
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 400,
		Body:       "Bad Request",
	}, nil
}
