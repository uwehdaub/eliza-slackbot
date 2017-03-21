package main

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/necrophonic/go-eliza"
	"github.com/nlopes/slack"
)

func getUserMention(username string) string {
	return "<@" + username + ">"
}

func reactOnChannel(event *slack.MessageEvent, user *slack.User) {
	userMention := getUserMention(user.Name)
	message := fmt.Sprintf("Hi %s! My name is Eliza.\n", userMention)
	message += "I believe you want this to be handled discreetly!\n"
	message += "So let us switch to private communication, please."

	log.Debug("Sending message: ", message)
	rtm.SendMessage(rtm.NewOutgoingMessage(message, event.Channel))
}

func reactOnGroup(event *slack.MessageEvent, user *slack.User) {
	userMention := getUserMention(user.Name)
	message := fmt.Sprintf("Hi %s! My name is Eliza.\n", userMention)
	message += "And I'm not a group therapist!\n"
	message += "Please let us switch to private communication."

	log.Debug("Sending message: ", message)
	rtm.SendMessage(rtm.NewOutgoingMessage(message, event.Channel))
}

func reactOnIm(event *slack.MessageEvent, user *slack.User) {
	userMention := getUserMention(user.Name)
	message := ""

	if isGreeting(event.Text) {
		message = fmt.Sprintf("Hello %s! My name is Eliza.\nHow can I help you?", userMention)
	} else if isBye(event.Text) {
		message = fmt.Sprintf("Goodbye %s! I hope to see you soon.", userMention)
	} else {
		message, _ = eliza.AnalyseString(event.Text)
	}

	log.Debug("Sending message: ", message)
	rtm.SendMessage(rtm.NewOutgoingMessage(message, event.Channel))
}

func isGreeting(s string) bool {
	switch {
	case strings.HasPrefix(s, botMention):
		return true
	case strings.HasPrefix(s, "Hi"):
		return true
	case strings.HasPrefix(s, "Hello"):
		return true
	case strings.Contains(strings.ToLower(s), "hello"):
		return true
	case strings.HasPrefix(s, "Good morning"):
		return true
	case strings.HasPrefix(s, "Good afternoon"):
		return true
	case strings.HasPrefix(s, "Good evening"):
		return true
	}
	return false
}

func isBye(s string) bool {
	switch {
	case strings.Contains(strings.ToLower(s), "bye"):
		return true
	case strings.HasPrefix(strings.ToLower(s), "goodbye"):
		return true
	case strings.HasPrefix(strings.ToLower(s), "exit"):
		return true
	case strings.HasPrefix(strings.ToLower(s), "quit"):
		return true
	}
	return false
}
