package main

import (
	"encoding/json"
	"log"
	"slack-lottery-bot/adaptor"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
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
			msgOption := slack.MsgOptionAttachments(postMsgParams(h.selectActionOptions()))
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

func (h *handler) selectActionOptions() []slack.AttachmentActionOption {
	options := []slack.AttachmentActionOption{
		{
			Text:  "このチャンネルのメンバーから",
			Value: "channel",
		},
	}

	// UserGroupから抽選するメニューを追加
	groups, err := h.api.GetUserGroups()
	if err != nil {
		log.Print(err)
		return nil
	}
	for _, group := range groups {
		options = append(options, slack.AttachmentActionOption{
			Text:  group.Name,
			Value: group.ID,
		})
	}
	return options
}

func postMsgParams(selectOptions []slack.AttachmentActionOption) slack.Attachment {
	attachment := slack.Attachment{
		Text:       "メニューを選んでください",
		Color:      "#f9a41b",
		CallbackID: "select",
		Actions: []slack.AttachmentAction{
			{
				Name:    "select",
				Type:    "select",
				Options: selectOptions,
			},
			{
				Name:  "cancel",
				Text:  "Cancel",
				Type:  "button",
				Style: "danger",
			},
		},
	}

	return attachment
}
