package main

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/url"
	"slack-lottery-bot/adaptor"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/slack-go/slack"
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

	if len(message.ActionCallback.AttachmentActions) == 0 {
		view := new(slack.ViewResponse)
		if err := json.Unmarshal([]byte(jsonStr), view); err != nil {
			log.Printf("failed to decode json message from slack: %s", jsonStr)
			return events.APIGatewayProxyResponse{}, err
		}

		values := view.State.Values
		target := values["target"]["target"].SelectedOption.Value
		count := values["count"]["count"].Value

		log.Println("count, terget, callbackID: ", count, target, view.CallbackID)
		err := h.lottery(target, count, view.CallbackID)
		code := 200
		if err != nil {
			code = 400
		}
		return events.APIGatewayProxyResponse{
			StatusCode: code,
		}, err
	}

	action := message.ActionCallback.AttachmentActions[0]
	switch action.Name {
	case "start":
		log.Print("start action")
		originalMessage := message.OriginalMessage
		jsonBody, err := responseMessage(&originalMessage, "抽選しました", "")
		if err != nil {
			log.Print(err)
			return events.APIGatewayProxyResponse{}, err
		}

		log.Print("select start TriggerID: ", message.TriggerID)
		err = h.api.OpenView(message.TriggerID, h.modalViewReqest(message.Channel.ID))
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

func (h *handler) modalViewReqest(channelID string) slack.ModalViewRequest {
	return slack.ModalViewRequest{
		Type:       slack.VTModal,
		CallbackID: channelID,
		Title:      slack.NewTextBlockObject("plain_text", "Lottery Bot", false, false),
		Close:      slack.NewTextBlockObject("plain_text", "Cancel", false, false),
		Submit:     slack.NewTextBlockObject("plain_text", "Submit", false, false),
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				slack.NewInputBlock(
					"target",
					slack.NewTextBlockObject("plain_text", "抽選対象", false, false),
					slack.NewOptionsSelectBlockElement(
						"static_select",
						slack.NewTextBlockObject("plain_text", "抽選対象", false, false),
						"target",
						h.selectElements(channelID)...,
					),
				),
				slack.NewInputBlock(
					"count",
					slack.NewTextBlockObject("plain_text", "抽選人数", false, false),
					slack.PlainTextInputBlockElement{
						Type:         "plain_text_input",
						ActionID:     "count",
						Placeholder:  slack.NewTextBlockObject("plain_text", "抽選人数", false, false),
						Multiline:    false,
						InitialValue: "1",
					},
				),
			},
		},
	}
}

func (h *handler) selectElements(channelID string) []*slack.OptionBlockObject {
	// UserGroupから抽選するメニューを追加
	groups, err := h.api.GetUserGroups()
	if err != nil {
		log.Print(err)
		return []*slack.OptionBlockObject{
			&slack.OptionBlockObject{
				Text:  slack.NewTextBlockObject("plain_text", "このチャンネルのメンバーから", false, false),
				Value: channelID,
			},
		}
	}

	options := make([]*slack.OptionBlockObject, 0, len(groups))

	options = append(options, &slack.OptionBlockObject{
		Text:  slack.NewTextBlockObject("plain_text", "このチャンネルのメンバーから", false, false),
		Value: channelID,
	})

	for _, group := range groups {
		options = append(options, &slack.OptionBlockObject{
			Text:  slack.NewTextBlockObject("plain_text", group.Name, false, false),
			Value: group.ID,
		})
	}
	return options
}

func (h *handler) lottery(actionValue string, countValue string, channelID string) error {
	var userIDs []string
	var err error

	log.Printf("value: %s", actionValue)
	userIDs, err = h.api.GetUsersFromChannel(actionValue)
	if err != nil {
		userIDs, err = h.api.GetUsersFromUserGroup(actionValue)
	}

	if err != nil {
		return err
	}

	c, err := strconv.Atoi(countValue)
	if err != nil {
		c = 1
	}

	lotteriedUids := lotteryUsersFromUsers(userIDs, c)
	us := make([]string, len(lotteriedUids))
	for _, uid := range lotteriedUids {
		us = append(us, "<@"+uid+">\n")
	}
	pMsg := strings.Join(us, "") + "が当選しました"
	return h.api.PostMessage(channelID, pMsg)
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

func lotteryUsersFromUsers(userIDs []string, count int) []string {
	rand.Seed(time.Now().UnixNano())
	l := len(userIDs)
	for i := l - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		userIDs[i], userIDs[j] = userIDs[j], userIDs[i]
	}
	return userIDs[:count]
}
