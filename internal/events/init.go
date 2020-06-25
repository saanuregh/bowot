package events

import "github.com/andersfylling/disgord"

var events = map[string]interface{}{}

func Register(client *disgord.Client) {
	for i, v := range events {
		client.On(i, v)
	}
}
