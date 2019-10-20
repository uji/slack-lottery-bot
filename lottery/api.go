package main

import (
	"math/rand"
	"time"

	"github.com/nlopes/slack"
)

func getOneUserFromChannel(channelID string) (string, error) {
	params := slack.GetUsersInConversationParameters{ChannelID: channelID}
	memberIDs, _, err := api.GetUsersInConversation(&params)
	if err != nil {
		return "", err
	}
	rand.Seed(time.Now().UnixNano())
	memberID := memberIDs[rand.Intn(len(memberIDs)-1)]
	// member, err := api.GetUserInfo(memberID)
	return memberID, err
}
