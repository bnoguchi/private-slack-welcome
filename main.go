package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"os"
)

var SLACK_TOKEN = os.Getenv("SLACK_TOKEN")
var SEND_AS_USER = os.Getenv("SEND_AS_USER")

func main() {
	api := slack.New(SLACK_TOKEN)
	rtm := api.NewRTM()

	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			fmt.Print("Event Received: ")
			switch ev := msg.Data.(type) {
			case *slack.TeamJoinEvent:
				welcomeUser(ev.User.ID, api)
			case *slack.LatencyReport:
				fmt.Printf("Current latency: %v\n", ev.Value)

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:

				// Ignore other events..
				// fmt.Printf("Unexpected: %v\n", msg.Data)
			}
		}
	}
}

func welcomeUser(user string, api *slack.Client) {
	fmt.Printf("New user joined: %s", user)
	_, _, channelID, err := api.OpenIMChannel(user)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer api.CloseIMChannel(channelID)
	text := `
Welcome to the Slack group! When you have a moment, head over to the #intros channel to introduce yourself to everyone here. We're making introducing yourself a requirement of being a part of the group. We want people to know they can reach out to a community of diverse experts on here to ask any questions or to engage in fruitful conversations. You should mention your name, what you do, and any expertise you have that might be helpful to the group. Thanks, and looking forward to having you here!
`
	params := slack.PostMessageParameters{
		Username:  SEND_AS_USER,
		AsUser:    true,
		LinkNames: 1,
	}
	_, timestamp, err := api.PostMessage(channelID, text, params)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}
