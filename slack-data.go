package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/nlopes/slack"
)

func getBotInfo(api *slack.Client) {
	info, err := api.AuthTest()
	if err != nil {
		log.Error("Error calling AuthTest: ", err)
		return
	}
	botUserID = info.UserID
	botMention = fmt.Sprintf("<@%s>", botUserID)
	log.Info("botUserID = ", botUserID)
}

func getChannelData(api *slack.Client) {
	channels, err := api.GetChannels(true)
	if err != nil {
		log.Error("Error getting Channels: ", err)
		return
	}
	// clear channelList ...
	channelList = map[string]slack.Channel{}
	// ... and reload
	log.Info("Adding to channelList:")
	for _, channel := range channels {
		channelList[channel.ID] = channel
		log.Info(fmt.Sprintf("Channel ID: %s / Name: %s", channel.ID, channel.Name))
	}
}

func getImData(api *slack.Client) {
	ims, err := api.GetIMChannels()
	if err != nil {
		log.Error("Error getting IMChannels: ", err)
		return
	}
	// clear imList ...
	imList = map[string]slack.IM{}
	// ... and reload
	log.Info("Adding to IMList:")
	for _, im := range ims {
		imList[im.ID] = im
		log.Info(fmt.Sprintf("IM ID: %s / User: %s", im.ID, im.User))
	}
}

func getGroupData(api *slack.Client) {
	groups, err := api.GetGroups(true)
	if err != nil {
		log.Error("Error getting Groups: ", err)
		return
	}
	// clear groupList ...
	groupList = map[string]slack.Group{}
	// ... and reload
	log.Info("Adding to GroupList:")
	for _, group := range groups {
		groupList[group.ID] = group
		log.Info(fmt.Sprintf("Group ID: %s / Name: %s / User: %+v", group.ID, group.Name, group.Members))
	}
}
