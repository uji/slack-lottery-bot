package adaptor

import (
	"log"

	"github.com/slack-go/slack"
)

type api struct {
	bot   *slack.Client
	oauth *slack.Client
}

type API interface {
	PostMessage(channel string, text string) error
	PostMessageWithOptions(change string, text string, msgOption slack.MsgOption) error
	OpenView(triggerID string, view slack.ModalViewRequest) error
	GetUsersFromChannel(channelID string) ([]string, error)
	GetUserGroups() ([]slack.UserGroup, error)
	GetUsersFromUserGroup(groupID string) ([]string, error)
}

func NewAPI(botToken string, oauthToken string) API {
	return &api{slack.New(botToken), slack.New(oauthToken)}
}

func (a *api) PostMessage(channel string, text string) error {
	_, _, err := a.bot.PostMessage(channel, slack.MsgOptionText(text, false))
	if err != nil {
		log.Println("PostMessageErr: ", err)
	}
	return err
}

func (a *api) PostMessageWithOptions(channel string, text string, msgOption slack.MsgOption) error {
	_, _, err := a.bot.PostMessage(channel, slack.MsgOptionText(text, false), msgOption)
	if err != nil {
		log.Println("PostMessageWithOptionsErr: ", err)
	}
	return err
}

func (a *api) OpenView(triggerID string, view slack.ModalViewRequest) error {
	_, err := a.bot.OpenView(triggerID, view)
	if err != nil {
		log.Println("OpenViewErr: ", err)
	}
	return err
}

func (a *api) GetUsersFromChannel(channelID string) ([]string, error) {
	params := slack.GetUsersInConversationParameters{ChannelID: channelID}
	userIDs, _, err := a.bot.GetUsersInConversation(&params)
	if err != nil {
		log.Println("GetUsersInConversationErr: ", err)
	}
	return userIDs, err
}

func (a *api) GetUserGroups() ([]slack.UserGroup, error) {
	groups, err := a.oauth.GetUserGroups()
	if err != nil {
		log.Println("GetUserGroupsErr: ", err)
	}
	return groups, err
}

func (a *api) GetUsersFromUserGroup(groupID string) ([]string, error) {
	userIDs, err := a.oauth.GetUserGroupMembers(groupID)
	if err != nil {
		log.Println("GetUsersFromUserGroupErr: ", err)
	}
	return userIDs, err
}
