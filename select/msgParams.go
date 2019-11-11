package main

import "github.com/nlopes/slack"

func postMsgParams(bot *slack.Client) slack.Attachment {
	attachment := slack.Attachment{
		Text:       "メニューを選んでください",
		Color:      "#f9a41b",
		CallbackID: "select",
		Actions: []slack.AttachmentAction{
			{
				Name:    "select",
				Type:    "select",
				Options: selectActionOptions(bot),
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

func selectActionOptions(bot *slack.Client) []slack.AttachmentActionOption {
	options := []slack.AttachmentActionOption{
		{
			Text:  "このチャンネルのメンバーから",
			Value: "channel",
		},
	}

	// UserGroupから抽選するメニューを追加
	groups, err := getUserGroups(bot)
	for _, group := range groups {
		options := append(options, slack.AttachmentActionOption{
			Text:  group.Name,
			Value: group.ID,
		})
	}
	return options
}
