package main

import (
	"encoding/json"
	"log"
	"slack-lottery-bot/adaptor"

	"github.com/aws/aws-lambda-go/events"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type handler struct {
	verificationToken string
	api               adaptor.API
}

type Handler interface {
	Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

func NewHandler(verificationToken string, botToken string, oauthToken string) Handler {
	return &handler{verificationToken, adaptor.NewAPI(botToken, oauthToken)}
}

func (h *handler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	reqBody := request.Body
	log.Printf("request body: %#v", reqBody)
	eventsAPIEvent, err := slackevents.ParseEvent(
		json.RawMessage(reqBody),
		slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: h.verificationToken}),
	)
	if err != nil {
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
			msgOption := slack.MsgOptionAttachments(postMsgParams())
			h.api.PostMessageWithOptions(ev.Channel, "ユーザーの抽選を始めます", msgOption)
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

func postMsgParams() slack.Attachment {
	return slack.Attachment{
		Text:       "モーダルで抽選方法を指定してください",
		Color:      "#f9a41b",
		CallbackID: "start",
		Actions: []slack.AttachmentAction{
			{
				Name: "start",
				Text: "Start!",
				Type: "button",
			},
			{
				Name:  "cancel",
				Text:  "Cancel",
				Type:  "button",
				Style: "danger",
			},
		},
	}
}
