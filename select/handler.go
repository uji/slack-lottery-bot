package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)

type handler struct {
	verificationToken string
	bot               *slack.Client
}

type Handler interface {
	Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

func NewHandler(verificationToken string, botToken string) Handler {
	return &handler{verificationToken, slack.New(botToken)}
}

func (h *handler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	reqBody := request.Body
	eventsAPIEvent, err := slackevents.ParseEvent(
		json.RawMessage(reqBody),
		slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: h.verificationToken}),
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
			h.bot.PostMessage(ev.Channel, slack.MsgOptionText("ユーザーの抽選を始めます", false), slack.MsgOptionAttachments(postMsgParams()))
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
