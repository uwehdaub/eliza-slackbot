package main

import (
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/nlopes/slack"
)

var (
	rtm *slack.RTM
	api *slack.Client

	channelList = map[string]slack.Channel{}
	imList      = map[string]slack.IM{}
	groupList   = map[string]slack.Group{}

	botUserID  = ""
	botMention = ""
)

func main() {

	log.SetLevel(log.DebugLevel)

	token := os.Getenv("SLACK_TOKEN")
	api = slack.New(token)
	rtm = api.NewRTM()
	go rtm.ManageConnection()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			log.Debug("Received event of type: ", msg.Type)

			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				log.Info("Processing event of type: ", msg.Type)
				log.Info("Connection counter: ", ev.ConnectionCount)

			case *slack.HelloEvent:
				log.Info("Processing event of type: ", msg.Type)
				getBotInfo(api)
				getChannelData(api)
				getImData(api)
				getGroupData(api)

			case *slack.ChannelCreatedEvent,
				*slack.ChannelArchiveEvent,
				*slack.ChannelRenameEvent:
				log.Info("Processing event of type: ", msg.Type)
				getChannelData(api)

			case *slack.IMCreatedEvent,
				*slack.UserChangeEvent:
				log.Info("Processing event of type: ", msg.Type)
				getImData(api)

			case *slack.GroupCreatedEvent,
				*slack.GroupRenameEvent,
				*slack.GroupCloseEvent,
				*slack.GroupJoinedEvent:
				log.Info("Processing event of type: ", msg.Type)
				getGroupData(api)

			case *slack.MessageEvent:
				log.Info("Processing event of type: ", msg.Type)
				log.Debug(fmt.Sprintf("Message: %+v\n", ev))

				info := rtm.GetInfo()

				if !ev.Hidden && !ownMessage(ev.User) {
					user := info.GetUserByID(ev.User)
					if user == nil || user.IsBot {
						log.Debug("No (real) user found.")
						break
					}

					switch {
					case isChannel(ev.Channel) && strings.Contains(ev.Text, botMention):
						go reactOnChannel(ev, user)
					case isIm(ev.Channel):
						go reactOnIm(ev, user)
					case isGroup(ev.Channel):
						go reactOnGroup(ev, user)
					}
				}

			case *slack.RTMError:
				log.Error(ev.Error())

			case *slack.InvalidAuthEvent:
				log.Fatal("Invalid credentials")

			default:
				//Take no action
			}
		}
	}
}

func ownMessage(UserID string) bool {
	return botUserID == UserID
}

func isChannel(channelID string) bool {
	_, found := channelList[channelID]
	return found
}
func isIm(imID string) bool {
	_, found := imList[imID]
	return found
}
func isGroup(groupID string) bool {
	_, found := groupList[groupID]
	return found
}
