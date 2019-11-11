package main

import (
	"math/rand"
	"time"

	"github.com/nlopes/slack"
)

func postResultMessage(bot *slack.Client, channel string, text string) error {
	_, _, err := bot.PostMessage(channel, slack.MsgOptionText(text, false))
	return err
}

func getUsersFromChannel(bot *slack.Client, channelID string) ([]string, error) {
	params := slack.GetUsersInConversationParameters{ChannelID: channelID}
	userIDs, _, err := bot.GetUsersInConversation(&params)
	return userIDs, err
}

func getUserGroups(bot *slack.Client) ([]slack.UserGroup, error) {
	groups, err := bot.GetUserGroups()
	return groups, err
}

func getUsersFromUserGroup(bot *slack.Client, groupID string) ([]string, error) {
	userIDs, err := bot.GetUserGroupMembers(groupID)
	return userIDs, err
}

func lotteryOneUserFromUsers(userIDs []string) string {
	rand.Seed(time.Now().UnixNano())
	userID := userIDs[rand.Intn(len(userIDs)-1)]
	return userID
}
